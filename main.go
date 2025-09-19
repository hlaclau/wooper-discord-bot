package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"wooper-bot/internal/bot"
	"wooper-bot/internal/config"
	"wooper-bot/internal/handlers"
	"wooper-bot/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	imageService, err := services.NewImageService("img")
	if err != nil {
		log.Fatalf("image service error: %v", err)
	}

	messageHandler := handlers.NewMessageHandler(imageService)

	b, err := bot.New(cfg.DiscordBotToken)
	if err != nil {
		log.Fatalf("bot error: %v", err)
	}
	b.AddHandler(messageHandler.OnMessageCreate)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := b.Start(ctx); err != nil {
		log.Fatalf("run error: %v", err)
	}
}
