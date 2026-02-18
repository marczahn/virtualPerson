package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/i18n"
	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/output"
	"github.com/marczahn/person/internal/psychology"
	"github.com/marczahn/person/internal/reviewer"
	"github.com/marczahn/person/internal/sense"
	"github.com/marczahn/person/internal/server"
	"github.com/marczahn/person/internal/simulation"
)

type config struct {
	AnthropicAPIKey string `json:"anthropic_api_key"`
	Model           string `json:"model"`
	DBPath          string `json:"db_path"`
	Lang            string `json:"lang"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig() config {
	cfgPath := "config.json"
	if p := os.Getenv("PERSON_CONFIG"); p != "" {
		cfgPath = p
	}

	var cfg config
	data, err := os.ReadFile(cfgPath)
	if err == nil {
		json.Unmarshal(data, &cfg)
	}
	return cfg
}

func run() error {
	serverMode    := flag.Bool("server", false, "run in server mode with WebSocket support")
	addr          := flag.String("addr", ":8080", "server listen address (server mode only)")
	lang          := flag.String("lang", "", "language code (e.g., en, de)")
	scenarioFile  := flag.String("scenario-file", "", "path to a text file describing the person's physical environment")
	flag.Parse()

	fileCfg := loadConfig()

	// Resolve language: flag > env > config > default.
	langCode := fileCfg.Lang
	if l := os.Getenv("PERSON_LANG"); l != "" {
		langCode = l
	}
	if *lang != "" {
		langCode = *lang
	}
	if langCode == "" {
		langCode = "en"
	}
	if err := i18n.Load(langCode); err != nil {
		return fmt.Errorf("load language %q: %w", langCode, err)
	}

	// Config file values, overridable by env vars.
	if key := os.Getenv("ANTHROPIC_API_KEY"); key == "" && fileCfg.AnthropicAPIKey != "" {
		os.Setenv("ANTHROPIC_API_KEY", fileCfg.AnthropicAPIKey)
	}

	dbPath := fileCfg.DBPath
	if p := os.Getenv("PERSON_DB"); p != "" {
		dbPath = p
	}
	if dbPath == "" {
		dbPath = "person.db"
	}

	modelStr := fileCfg.Model
	if m := os.Getenv("PERSON_MODEL"); m != "" {
		modelStr = m
	}
	model := anthropic.Model(modelStr)
	if model == "" {
		model = anthropic.ModelClaudeHaiku4_5
	}

	store, err := memory.NewSQLiteStore(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer store.Close()

	// Load or create biological state.
	bioState, err := store.LoadBioState()
	if err != nil {
		return fmt.Errorf("load bio state: %w", err)
	}
	tr := i18n.T()
	if bioState == nil {
		s := biology.NewDefaultState()
		bioState = &s
		fmt.Println(tr.CLI.FreshState)
	} else {
		bioState.LastUpdate = time.Now()
		fmt.Println(tr.CLI.ResumedState)
	}

	// Load or create personality.
	personality, err := store.LoadPersonality()
	if err != nil {
		return fmt.Errorf("load personality: %w", err)
	}
	if personality == nil {
		personality = &psychology.Personality{
			Openness:          0.6,
			Conscientiousness: 0.5,
			Extraversion:      0.4,
			Agreeableness:     0.6,
			Neuroticism:       0.5,
		}
		if err := store.SavePersonality(personality); err != nil {
			return fmt.Errorf("save personality: %w", err)
		}
		fmt.Println(tr.CLI.CreatedPersonality)
	}

	// Load or create identity core.
	identity, err := store.LoadIdentityCore()
	if err != nil {
		return fmt.Errorf("load identity: %w", err)
	}
	if identity == nil {
		identity = &memory.IdentityCore{
			SelfNarrative:     tr.Defaults.SelfNarrative,
			DispositionTraits: tr.Defaults.DispositionTraits,
			RelationalMarkers: []string{},
			KeyMemories:       []string{},
			EmotionalPatterns: tr.Defaults.EmotionalPatterns,
			ValuesCommitments: tr.Defaults.ValuesCommitments,
			LastUpdated:       time.Now(),
		}
		if err := store.SaveIdentityCore(identity); err != nil {
			return fmt.Errorf("save identity: %w", err)
		}
		fmt.Println(tr.CLI.CreatedIdentity)
	}

	// Load episodic memories for consciousness context.
	memories, err := store.LoadMemories()
	if err != nil {
		return fmt.Errorf("load memories: %w", err)
	}

	// Load scenario file if provided.
	var scenario string
	if *scenarioFile != "" {
		data, err := os.ReadFile(*scenarioFile)
		if err != nil {
			return fmt.Errorf("read scenario file: %w", err)
		}
		scenario = strings.TrimSpace(string(data))
		fmt.Printf(tr.CLI.ScenarioLoaded+"\n", *scenarioFile)
	}

	// Build components.
	llm := consciousness.NewClaudeAdapter(model)

	consciousnessEngine := consciousness.NewEngine(consciousness.EngineConfig{
		LLM:                 llm,
		Identity:            identity,
		MaxPromptTokens:     2000,
		MaxContextMemories:  5,
		MinCallInterval:     2 * time.Second,
		SpontaneousInterval: 30 * time.Second,
	})
	consciousnessEngine.UpdateMemories(memories)

	display := output.NewDisplay(os.Stdout, true)

	psychReviewer := reviewer.NewReviewer(reviewer.ReviewerConfig{
		LLM:         llm,
		MinInterval: 60 * time.Second,
		MaxThoughts: 20,
	})

	// Determine input source: stdin (default) or pipe from WebSocket hub.
	var input io.Reader = os.Stdin
	var hub *server.Hub
	var pipeWriter *io.PipeWriter

	if *serverMode {
		pr, pw := io.Pipe()
		pipeWriter = pw
		input = pr
		hub = server.NewHub(pw)

		// Broadcast MIND entries to connected WebSocket clients.
		display.SetListener(func(entry output.Entry) {
			if entry.Source != output.Mind {
				return
			}
			hub.Broadcast(server.ServerMessage{
				Type:        "thought",
				Content:     entry.Message,
				ThoughtType: entry.ThoughtType,
				Trigger:     entry.Trigger,
				Timestamp:   entry.Timestamp,
			})
		})
	}

	cfg := simulation.Config{
		BioProcessor:   biology.NewProcessor(),
		PsychProcessor: psychology.NewProcessor(*personality),
		Consciousness:  consciousnessEngine,
		SenseParser:    sense.NewKeywordParser(),
		Display:        display,
		Store:          store,
		Reviewer:       psychReviewer,
		Personality:    personality,
		BioState:       bioState,
		Identity:       identity,
		Scenario:       scenario,
		TickInterval:   100 * time.Millisecond,
		SimStart:       time.Date(2024, 6, 15, 8, 0, 0, 0, time.Local),
		Input:          input,
	}

	loop := simulation.NewLoop(cfg)

	// Handle signals for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println(tr.CLI.ShuttingDown)
		cancel()
	}()

	if *serverMode {
		httpSrv := &http.Server{
			Addr:    *addr,
			Handler: server.NewHandler(hub),
		}
		go func() {
			fmt.Printf(tr.CLI.WebSocketListening+"\n", *addr)
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Fprintf(os.Stderr, "http server error: %v\n", err)
			}
		}()
		defer func() {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			httpSrv.Shutdown(shutdownCtx)
			pipeWriter.Close()
		}()
	}

	fmt.Println(tr.CLI.SimulationStarted)
	fmt.Printf(tr.CLI.UsingModel+"\n", model)

	return loop.Run(ctx)
}
