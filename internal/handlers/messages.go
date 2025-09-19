package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"wooper-bot/internal/logger"
	"wooper-bot/internal/services"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type MessageHandler struct {
	ImageService *services.ImageService
}

func NewMessageHandler(imageService *services.ImageService) *MessageHandler {
	return &MessageHandler{ImageService: imageService}
}

func (h *MessageHandler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}

	content := strings.TrimSpace(m.Content)

	// Log all messages for debugging (can be filtered by log level)
	logger.Logger.Debug("Message received",
		zap.String("user", m.Author.Username),
		zap.String("user_id", m.Author.ID),
		zap.String("channel_id", m.ChannelID),
		zap.String("guild_id", m.GuildID),
		zap.String("content", content))

	// Check if message starts with ! and has a valid category
	if strings.HasPrefix(content, "!") {
		category := strings.TrimPrefix(content, "!")

		// Log command attempt
		logger.Logger.Info("Command received",
			zap.String("command", content),
			zap.String("category", category),
			zap.String("user", m.Author.Username),
			zap.String("user_id", m.Author.ID),
			zap.String("channel_id", m.ChannelID),
			zap.String("guild_id", m.GuildID))

		if h.ImageService.HasCategory(category) {
			startTime := time.Now()

			imagePath := h.ImageService.GetRandomImage(category)
			if imagePath == "" {
				logger.Logger.Warn("No images available for category",
					zap.String("category", category),
					zap.String("user", m.Author.Username))
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("no %s images available", category))
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			reader, fileName, err := h.ImageService.GetImageFile(ctx, imagePath)
			if err != nil {
				logger.Logger.Error("Failed to load image file",
					zap.String("category", category),
					zap.String("image_path", imagePath),
					zap.String("user", m.Author.Username),
					zap.Error(err))
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("failed to load %s: %v", category, err))
				return
			}
			defer reader.Close()

			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Files: []*discordgo.File{{
				Name:   fileName,
				Reader: reader,
			}}})

			duration := time.Since(startTime)

			if err != nil {
				logger.Logger.Error("Failed to send image",
					zap.String("category", category),
					zap.String("filename", fileName),
					zap.String("user", m.Author.Username),
					zap.Duration("duration", duration),
					zap.Error(err))
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("failed to send %s: %v", category, err))
			} else {
				logger.Logger.Info("Image sent successfully",
					zap.String("category", category),
					zap.String("filename", fileName),
					zap.String("user", m.Author.Username),
					zap.String("user_id", m.Author.ID),
					zap.String("channel_id", m.ChannelID),
					zap.Duration("duration", duration))
			}
		} else if category == "help" || category == "list" {
			// Show available categories
			logger.Logger.Info("Help command requested",
				zap.String("user", m.Author.Username),
				zap.String("user_id", m.Author.ID))

			categories := h.ImageService.GetAvailableCategories()
			if len(categories) == 0 {
				logger.Logger.Warn("No categories available for help",
					zap.String("user", m.Author.Username))
				_, _ = s.ChannelMessageSend(m.ChannelID, "no image categories available")
				return
			}

			message := "Available image categories:\n"
			for _, cat := range categories {
				count := h.ImageService.GetImageCount(cat)
				message += fmt.Sprintf("• `!%s` (%d images)\n", cat, count)
			}

			logger.Logger.Info("Help response sent",
				zap.String("user", m.Author.Username),
				zap.Int("categories_count", len(categories)))

			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		} else {
			// Unknown command
			logger.Logger.Info("Unknown command received",
				zap.String("command", content),
				zap.String("category", category),
				zap.String("user", m.Author.Username),
				zap.String("user_id", m.Author.ID))
		}
	}
}
