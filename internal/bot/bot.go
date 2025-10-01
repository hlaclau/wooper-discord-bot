package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func New(token string) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent
	return &Bot{session: dg}, nil
}

func (b *Bot) AddHandler(handler interface{}) func() {
	return b.session.AddHandler(handler)
}

func (b *Bot) RegisterSlashCommands(commands []*discordgo.ApplicationCommand) error {
	// Wait for the session to be ready
	if b.session.State.User == nil {
		return fmt.Errorf("session not ready, user is nil")
	}

	for _, cmd := range commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("create application command %s: %w", cmd.Name, err)
		}
		log.Printf("Registered slash command: /%s", cmd.Name)
	}
	return nil
}

func (b *Bot) Start(ctx context.Context) error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("open discord session: %w", err)
	}
	<-ctx.Done()
	return b.session.Close()
}

func (b *Bot) StartWithCommands(ctx context.Context, commands []*discordgo.ApplicationCommand) error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("open discord session: %w", err)
	}

	// Register slash commands after session is open
	if err := b.RegisterSlashCommands(commands); err != nil {
		return fmt.Errorf("register slash commands: %w", err)
	}

	<-ctx.Done()
	return b.session.Close()
}
