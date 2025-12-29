package handlers

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"viitorbot/internal/wiki"

	"github.com/bwmarrin/discordgo"
	"github.com/kubaliski/golog/pkg/logger"
)

type fakeWiki struct{}

func (f *fakeWiki) GetRandomArticle() (string, string, string, error) {
	return "T", "E", "https://example.org/t", nil
}

func (f *fakeWiki) GetRandomArticleByDate(t time.Time) (string, string, string, *wiki.Evidence, error) {
	return "DT", "DE", "https://example.org/dt", &wiki.Evidence{Type: "event", Year: 2020, Text: "Evento de prueba"}, nil
}

type fakeWikiErr struct{}

func (f *fakeWikiErr) GetRandomArticle() (string, string, string, error) {
	return "", "", "", errors.New("upstream")
}

func (f *fakeWikiErr) GetRandomArticleByDate(t time.Time) (string, string, string, *wiki.Evidence, error) {
	return "", "", "", nil, errors.New("upstream")
}

func TestBuildResponse(t *testing.T) {
	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	h := New(lg, &fakeWiki{}, base)

	resp, err := h.BuildResponse(base, "Hola viitorbot hazlotuyo por favor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == "" {
		t.Fatalf("expected a response but got empty string")
	}
	if !strings.Contains(resp, "T") || !strings.Contains(resp, "https://example.org/t") {
		t.Fatalf("response missing expected parts: %q", resp)
	}

	// message without trigger
	resp2, err := h.BuildResponse(base, "nothing to see here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp2 != "" {
		t.Fatalf("expected empty response when not triggered, got %q", resp2)
	}
}

// --- HandleMessage tests using a fake SessionAdapter ---

type fakeAdapter struct {
	botID       string
	lastChannel string
	lastContent string
	sendErr     error
}

func (f *fakeAdapter) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	f.lastChannel = channelID
	f.lastContent = content
	if f.sendErr != nil {
		return nil, f.sendErr
	}
	return &discordgo.Message{ID: "1"}, nil
}

func (f *fakeAdapter) BotUserID() string { return f.botID }

func TestHandleMessage_SendsResponse(t *testing.T) {
	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	h := New(lg, &fakeWiki{}, base)

	fa := &fakeAdapter{botID: "bot123"}
	m := &discordgo.MessageCreate{Message: &discordgo.Message{Author: &discordgo.User{ID: "user1"}, ChannelID: "chan1", Content: "viitorbot hazlotuyo"}}

	h.HandleMessage(fa, m)

	if fa.lastChannel != "chan1" {
		t.Fatalf("expected message sent to chan1, got %q", fa.lastChannel)
	}
	if fa.lastContent == "" {
		t.Fatalf("expected content to be sent")
	}
}

func TestHandleMessage_WikiError_NotifiesChannel(t *testing.T) {
	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	h := New(lg, &fakeWikiErr{}, base)

	fa := &fakeAdapter{botID: "bot123"}
	m := &discordgo.MessageCreate{Message: &discordgo.Message{Author: &discordgo.User{ID: "user1"}, ChannelID: "chanX", Content: "viitorbot hazlotuyo"}}

	h.HandleMessage(fa, m)

	if fa.lastChannel != "chanX" {
		t.Fatalf("expected notification sent to chanX, got %q", fa.lastChannel)
	}
	if fa.lastContent == "" || !strings.Contains(fa.lastContent, "Error procesando mensaje") {
		t.Fatalf("expected error notification content, got %q", fa.lastContent)
	}
}

func TestHandleMessage_IgnoresBotAuthor(t *testing.T) {
	base := context.Background()
	base = logger.SetServiceName(base, "test")
	lg := logger.FromContext(base)

	h := New(lg, &fakeWiki{}, base)

	fa := &fakeAdapter{botID: "bot123"}
	m := &discordgo.MessageCreate{Message: &discordgo.Message{Author: &discordgo.User{ID: "bot123"}, ChannelID: "chanY", Content: "viitorbot hazlotuyo"}}

	h.HandleMessage(fa, m)

	if fa.lastContent != "" {
		t.Fatalf("expected no message sent when author is bot, got %q", fa.lastContent)
	}
}
