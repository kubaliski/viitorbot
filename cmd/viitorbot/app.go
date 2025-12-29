package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kubaliski/golog/pkg/logger"
)

// BotLifecycle defines the minimal lifecycle methods required by Run.
type BotLifecycle interface {
	Start() error
	Stop() error
}

// Run starts the provided bot (created by createBot), waits for a signal on stopCh,
// and stops the bot. Returns any error encountered during creation, start, or stop.
func Run(createBot func() (BotLifecycle, error), stopCh <-chan os.Signal, lg *logger.Logger) error {
	b, err := createBot()
	if err != nil {
		if lg != nil {
			lg.Error("error creating bot: " + err.Error())
		}
		return err
	}

	if err := b.Start(); err != nil {
		if lg != nil {
			lg.Error("error starting bot: " + err.Error())
		}
		return err
	}

	// Wait for stop signal
	<-stopCh

	if err := b.Stop(); err != nil {
		if lg != nil {
			lg.Error("error stopping bot: " + err.Error())
		}
		return err
	}

	return nil
}

// DefaultStopChannel returns a channel that receives OS interrupt/terminate signals.
func DefaultStopChannel() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}
