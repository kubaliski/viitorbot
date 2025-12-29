package bot

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/kubaliski/golog/pkg/logger"

	"viitorbot/internal/config"
	"viitorbot/internal/handlers"
)

type fakeSession struct {
	addHandlerCalls int
	openCalled      bool
	closeCalled     bool
	openErr         error
	closeErr        error
	handlers        []interface{}
}

func (f *fakeSession) AddHandler(handler interface{}) func() {
	f.addHandlerCalls++
	f.handlers = append(f.handlers, handler)
	return func() {}
}
func (f *fakeSession) Open() error  { f.openCalled = true; return f.openErr }
func (f *fakeSession) Close() error { f.closeCalled = true; return f.closeErr }

func TestNewStartStop(t *testing.T) {
	old := defaultSessionFactory
	fs := &fakeSession{}
	defaultSessionFactory = func(token string) (sessionIface, error) { return fs, nil }
	defer func() { defaultSessionFactory = old }()

	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	cfg := &config.Config{DiscordToken: "token123"}
	h := handlers.New(lg, nil, base)

	b, err := New(cfg, lg, base, h)
	if err != nil {
		t.Fatalf("unexpected error creating bot: %v", err)
	}

	if fs.addHandlerCalls < 2 {
		t.Fatalf("expected handlers registered, got %d", fs.addHandlerCalls)
	}

	// Invoke any ready handler to cover its body.
	if len(fs.handlers) > 0 {
		for _, hh := range fs.handlers {
			if hf, ok := hh.(func(*discordgo.Session, *discordgo.Ready)); ok {
				s := &discordgo.Session{}
				hf(s, &discordgo.Ready{})
			}
		}
	}

	if err := b.Start(); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}
	if !fs.openCalled {
		t.Fatalf("expected Open to be called")
	}

	if err := b.Stop(); err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}
	if !fs.closeCalled {
		t.Fatalf("expected Close to be called")
	}
}

func TestNew_FactoryError(t *testing.T) {
	old := defaultSessionFactory
	defaultSessionFactory = func(token string) (sessionIface, error) { return nil, fmt.Errorf("bad") }
	defer func() { defaultSessionFactory = old }()

	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	cfg := &config.Config{DiscordToken: "token123"}
	h := handlers.New(lg, nil, base)

	if _, err := New(cfg, lg, base, h); err == nil {
		t.Fatalf("expected error when session factory fails")
	}
}

func TestStart_Stop_Errors(t *testing.T) {
	old := defaultSessionFactory
	fs := &fakeSession{openErr: fmt.Errorf("openfail"), closeErr: fmt.Errorf("closefail")}
	defaultSessionFactory = func(token string) (sessionIface, error) { return fs, nil }
	defer func() { defaultSessionFactory = old }()

	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	cfg := &config.Config{DiscordToken: "token123"}
	h := handlers.New(lg, nil, base)

	b, err := New(cfg, lg, base, h)
	if err != nil {
		t.Fatalf("unexpected error creating bot: %v", err)
	}

	if err := b.Start(); err == nil {
		t.Fatalf("expected start error due to openErr")
	}
	if !fs.openCalled {
		t.Fatalf("expected Open to be called")
	}

	if err := b.Stop(); err == nil {
		t.Fatalf("expected stop error due to closeErr")
	}
	if !fs.closeCalled {
		t.Fatalf("expected Close to be called")
	}
}

func TestDefaultSessionFactory_CreatesSession(t *testing.T) {
	s, err := defaultSessionFactory("Bot dummy")
	if err != nil {
		t.Fatalf("expected no error from defaultSessionFactory, got %v", err)
	}
	if s == nil {
		t.Fatalf("expected session, got nil")
	}
	// ensure AddHandler exists and can be called
	s.AddHandler(func(*discordgo.Session, *discordgo.Ready) {})
}
