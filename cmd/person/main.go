package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/marczahn/person/internal/biology"
	"github.com/marczahn/person/internal/consciousness"
	"github.com/marczahn/person/internal/memory"
	"github.com/marczahn/person/internal/output"
	"github.com/marczahn/person/internal/psychology"
	"github.com/marczahn/person/internal/reviewer"
	"github.com/marczahn/person/internal/sense"
	"github.com/marczahn/person/internal/simulation"
)

type config struct {
	AnthropicAPIKey string `json:"anthropic_api_key"`
	Model           string `json:"model"`
	DBPath          string `json:"db_path"`
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
	fileCfg := loadConfig()

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
	if bioState == nil {
		s := biology.NewDefaultState()
		bioState = &s
		fmt.Println("Starting with fresh biological state.")
	} else {
		bioState.LastUpdate = time.Now()
		fmt.Println("Resumed from saved biological state.")
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
		fmt.Println("Created default personality.")
	}

	// Load or create identity core.
	identity, err := store.LoadIdentityCore()
	if err != nil {
		return fmt.Errorf("load identity: %w", err)
	}
	if identity == nil {
		identity = &memory.IdentityCore{
			SelfNarrative:     "I'm a person, trying to make sense of the world around me. I think a lot, sometimes too much.",
			DispositionTraits: []string{"thoughtful", "curious", "sometimes anxious"},
			RelationalMarkers: []string{},
			KeyMemories:       []string{},
			EmotionalPatterns: []string{"tends to overthink under stress", "finds comfort in routines"},
			ValuesCommitments: []string{"honesty", "understanding", "being kind"},
			LastUpdated:       time.Now(),
		}
		if err := store.SaveIdentityCore(identity); err != nil {
			return fmt.Errorf("save identity: %w", err)
		}
		fmt.Println("Created default identity.")
	}

	// Load episodic memories for consciousness context.
	memories, err := store.LoadMemories()
	if err != nil {
		return fmt.Errorf("load memories: %w", err)
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
		TickInterval:   100 * time.Millisecond,
		SimStart:       time.Date(2024, 6, 15, 8, 0, 0, 0, time.Local),
		Input:          os.Stdin,
	}

	loop := simulation.NewLoop(cfg)

	// Handle signals for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	fmt.Println("Simulation started. Type to interact, Ctrl+C to quit.")
	fmt.Printf("Using model: %s\n\n", model)

	return loop.Run(ctx)
}
