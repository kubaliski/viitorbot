package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kubaliski/golog/pkg/logger"

	"viitorbot/internal/wiki"
)

type Handlers struct {
	lg      *logger.Logger
	wiki    wiki.Client
	baseCtx context.Context
}

func New(lg *logger.Logger, w wiki.Client, baseCtx context.Context) *Handlers {
	return &Handlers{lg: lg, wiki: w, baseCtx: baseCtx}
}

func (h *Handlers) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Adapt to new HandleMessage using a real adapter so MessageCreate remains small.
	adapter := &realSessionAdapter{s: s}
	h.HandleMessage(adapter, m)
}

// SessionAdapter abstracts the minimal session behavior used by handlers, making it
// easy to test without a real discordgo.Session.
type SessionAdapter interface {
	ChannelMessageSend(channelID, content string) (*discordgo.Message, error)
	BotUserID() string
}

type realSessionAdapter struct {
	s *discordgo.Session
}

func (r *realSessionAdapter) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	return r.s.ChannelMessageSend(channelID, content)
}

func (r *realSessionAdapter) BotUserID() string { return r.s.State.User.ID }

// HandleMessage processes a message using a SessionAdapter (testable).
func (h *Handlers) HandleMessage(s SessionAdapter, m *discordgo.MessageCreate) {
	if m.Author.ID == s.BotUserID() {
		return
	}

	msgCtx := logger.WithServiceName(h.baseCtx, "message-create")

	resp, err := h.BuildResponse(msgCtx, m.Content)
	if err != nil {
		h.lg.Log(msgCtx, "ERROR", fmt.Sprintf("handler error: %v", err))
		// try to notify channel; ignore send errors
		s.ChannelMessageSend(m.ChannelID, "Error procesando mensaje: "+err.Error())
		return
	}

	if resp == "" {
		return
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, resp); err != nil {
		h.lg.Log(msgCtx, "ERROR", fmt.Sprintf("Error al enviar mensaje a Discord: %v", err))
	} else {
		h.lg.Log(msgCtx, "INFO", fmt.Sprintf("Respuesta enviada al canal %s (msg id: %s)", m.ChannelID, m.ID))
	}
}

// BuildResponse composes the reply for a given message content. Returns empty string when
// message does not trigger a response. This is separated for easier testing.
func (h *Handlers) BuildResponse(ctx context.Context, content string) (string, error) {
	lc := strings.ToLower(content)

	// detect date keywords/formats
	if dt, ok := parseDateFromContent(lc); ok {
		// Log the detected date for observability/debugging
		if h.lg != nil {
			h.lg.Log(ctx, "INFO", fmt.Sprintf("Detected date %s from content: %q", dt.Format("2006-01-02"), content))
		}
		title, extract, url, evidence, err := h.wiki.GetRandomArticleByDate(dt)
		if err != nil {
			// fallback to random article
			title, extract, url, err = h.wiki.GetRandomArticle()
			if err != nil {
				return "", fmt.Errorf("error obtener artículo: %w", err)
			}
		} else {
			// include evidence when available
			if evidence != nil {
				// log the evidence for debugging
				if h.lg != nil {
					h.lg.Log(ctx, "INFO", fmt.Sprintf("Evidence for selected item: type=%s year=%d text=%s", evidence.Type, evidence.Year, evidence.Text))
				}
			}
		}
		response := fmt.Sprintf("**%s**\n\n%s\n\nMás información: %s", title, extract, url)
		if evidence != nil {
			// append short provenance
			response = response + fmt.Sprintf("\n\nEvidencia: %s (%d) — %s", evidence.Type, evidence.Year, evidence.Text)
		}
		return response, nil
	}

	if strings.Contains(lc, "viitorbot hazlotuyo") {
		title, extract, url, err := h.wiki.GetRandomArticle()
		if err != nil {
			return "", fmt.Errorf("error obtener artículo: %w", err)
		}
		response := fmt.Sprintf("**%s**\n\n%s\n\nMás información: %s", title, extract, url)
		return response, nil
	}
	return "", nil
}

// parseDateFromContent attempts to find a date in common formats; returns the parsed time in local timezone.
func parseDateFromContent(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	// keywords
	if strings.Contains(s, "hoy") {
		return time.Now(), true
	}
	if strings.Contains(s, "ayer") {
		return time.Now().AddDate(0, 0, -1), true
	}
	if strings.Contains(s, "mañana") || strings.Contains(s, "manana") {
		return time.Now().AddDate(0, 0, 1), true
	}

	// try YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, true
	}
	// try DD/MM/YYYY
	if t, err := time.Parse("02/01/2006", s); err == nil {
		return t, true
	}
	// try DD-MM-YYYY
	if t, err := time.Parse("02-01-2006", s); err == nil {
		return t, true
	}

	// Not found
	return time.Time{}, false
}
