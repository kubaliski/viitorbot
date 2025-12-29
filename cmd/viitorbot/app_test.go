package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/kubaliski/golog/pkg/logger"
)

type fakeBot struct {
	startErr error
	stopErr  error
	started  bool
	stopped  bool
}

func (f *fakeBot) Start() error { f.started = true; return f.startErr }
func (f *fakeBot) Stop() error  { f.stopped = true; return f.stopErr }

func TestRun_SuccessLifecycle(t *testing.T) {
	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	fb := &fakeBot{}
	create := func() (BotLifecycle, error) { return fb, nil }

	stopCh := make(chan os.Signal, 1)

	done := make(chan error, 1)
	go func() {
		done <- Run(create, stopCh, lg)
	}()

	// give Run time to start
	time.Sleep(10 * time.Millisecond)
	if !fb.started {
		t.Fatalf("expected bot started")
	}

	stopCh <- os.Interrupt
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("Run did not return after stop signal")
	}
	if !fb.stopped {
		t.Fatalf("expected bot stopped")
	}
}

func TestRun_CreateError(t *testing.T) {
	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	create := func() (BotLifecycle, error) { return nil, errors.New("fail create") }
	stopCh := make(chan os.Signal, 1)

	if err := Run(create, stopCh, lg); err == nil {
		t.Fatalf("expected error from Run when create fails")
	}
}

func TestRun_StartError(t *testing.T) {
	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	fb := &fakeBot{startErr: errors.New("start fail")}
	create := func() (BotLifecycle, error) { return fb, nil }
	stopCh := make(chan os.Signal, 1)

	if err := Run(create, stopCh, lg); err == nil {
		t.Fatalf("expected error from Run when start fails")
	}
}
