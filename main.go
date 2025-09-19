package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"wooper-bot/internal/bot"
	"wooper-bot/internal/config"
	"wooper-bot/internal/handlers"
	"wooper-bot/internal/logger"
	"wooper-bot/internal/services"

	"go.uber.org/zap"
)

func main() {
	// Initialize logging
	if err := logger.Init(); err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	defer logger.Close()

	logger.Logger.Info("Starting wooper-bot")

	cfg, err := config.Load()
	if err != nil {
		logger.Logger.Fatal("config error", zap.Error(err))
	}

	imageService, err := services.NewImageService("img")
	if err != nil {
		logger.Logger.Fatal("image service error", zap.Error(err))
	}

	messageHandler := handlers.NewMessageHandler(imageService)

	b, err := bot.New(cfg.DiscordBotToken)
	if err != nil {
		logger.Logger.Fatal("bot error", zap.Error(err))
	}
	b.AddHandler(messageHandler.OnMessageCreate)

	logger.Logger.Info("Bot initialized successfully")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := b.Start(ctx); err != nil {
		logger.Logger.Fatal("run error", zap.Error(err))
	}
}
