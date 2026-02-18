package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marczahn/person/internal/client"
	"github.com/marczahn/person/internal/i18n"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	serverURL := flag.String("server", "ws://localhost:8080/ws", "WebSocket server URL")
	lang := flag.String("lang", "", "language code (e.g., en, de)")
	flag.Parse()

	// Resolve language: flag > env > default.
	langCode := "en"
	if l := os.Getenv("PERSON_LANG"); l != "" {
		langCode = l
	}
	if *lang != "" {
		langCode = *lang
	}
	if err := i18n.Load(langCode); err != nil {
		return fmt.Errorf("load language %q: %w", langCode, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	conn, err := client.Dial(ctx, *serverURL)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", *serverURL, err)
	}
	defer conn.Close()

	go conn.Run(ctx)

	model := client.NewModel(ctx, conn)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}
	return nil
}
