package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kubaliski/golog/pkg/logger"

	"viitorbot/internal/config"
	"viitorbot/internal/handlers"
)

// sessionIface abstracts the subset of discordgo.Session used by the bot.
type sessionIface interface {
	AddHandler(handler interface{}) func()
	Open() error
	Close() error
}

// defaultSessionFactory creates a real discord session. Tests can override this.
var defaultSessionFactory = func(token string) (sessionIface, error) {
	s, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type Bot struct {
	session  sessionIface
	handlers *handlers.Handlers
	lg       *logger.Logger
	ctx      context.Context
}

func New(cfg *config.Config, lg *logger.Logger, baseCtx context.Context, h *handlers.Handlers) (*Bot, error) {
	sess, err := defaultSessionFactory("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("error creating discord session: %w", err)
	}

	b := &Bot{session: sess, handlers: h, lg: lg, ctx: baseCtx}

	// Register ready handler and message handler
	// Use the underlying discordgo.Session when calling handlers that expect that type.
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		readyCtx := logger.WithServiceName(baseCtx, "ready-handler")
		user := "<unknown>"
		disc := "<unknown>"
		if s != nil && s.State != nil && s.State.User != nil {
			if s.State.User.Username != "" {
				user = s.State.User.Username
			}
			if s.State.User.Discriminator != "" {
				disc = s.State.User.Discriminator
			}
		}
		lg.Log(readyCtx, "INFO", fmt.Sprintf("Bot conectado como: %v#%v", user, disc))
	})

	b.session.AddHandler(h.MessageCreate)

	return b, nil
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("error opening discord session: %w", err)
	}
	b.lg.Info("Discord session opened")
	return nil
}

func (b *Bot) Stop() error {
	if err := b.session.Close(); err != nil {
		return fmt.Errorf("error closing discord session: %w", err)
	}
	b.lg.Info("Discord session closed")
	return nil
}
