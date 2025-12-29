package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kubaliski/golog/pkg/logger"

	"viitorbot/internal/bot"
	"viitorbot/internal/config"
	"viitorbot/internal/handlers"
	"viitorbot/internal/wiki"
)

func main() {
	baseCtx := context.Background()
	baseCtx = logger.SetServiceName(baseCtx, "viitorbot")
	lg := logger.FromContext(baseCtx)

	lg.Info("Iniciando bot...")

	// Use Run with the default create function and system stop channel.
	stopCh := DefaultStopChannel()
	if err := Run(createDefaultBot, stopCh, lg); err != nil {
		lg.Error(fmt.Sprintf("Error en ejecución: %v", err))
		os.Exit(1)
	}
}

// createDefaultBot constructs the real bot used by main. It's a variable so tests
// can override it when needed.
var createDefaultBot = func() (BotLifecycle, error) {
	baseCtx := context.Background()
	baseCtx = logger.SetServiceName(baseCtx, "viitorbot")
	lg := logger.FromContext(baseCtx)

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	wikiClient := wiki.NewClient()
	hs := handlers.New(lg, wikiClient, baseCtx)

	return bot.New(cfg, lg, baseCtx, hs)
}
