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

type InteractionHandler struct {
	ImageService *services.ImageService
}

func NewInteractionHandler(imageService *services.ImageService) *InteractionHandler {
	return &InteractionHandler{ImageService: imageService}
}

func (h *InteractionHandler) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "image" {
		h.handleImageCommand(s, i)
	}
}

func (h *InteractionHandler) handleImageCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	startTime := time.Now()

	// Get the category from the command options
	var category string
	if len(i.ApplicationCommandData().Options) > 0 {
		category = i.ApplicationCommandData().Options[0].StringValue()
	}

	// Log the interaction
	logger.Logger.Info("Slash command received",
		zap.String("command", "image"),
		zap.String("category", category),
		zap.String("user", i.Member.User.Username),
		zap.String("user_id", i.Member.User.ID),
		zap.String("channel_id", i.ChannelID),
		zap.String("guild_id", i.GuildID))

	// Check if category exists
	if !h.ImageService.HasCategory(category) {
		availableCategories := h.ImageService.GetAvailableCategories()
		message := fmt.Sprintf("Category '%s' not found. Available categories: %s",
			category, strings.Join(availableCategories, ", "))

		logger.Logger.Warn("Invalid category requested",
			zap.String("category", category),
			zap.Strings("available_categories", availableCategories),
			zap.String("user", i.Member.User.Username))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
		return
	}

	// Get random image
	imagePath := h.ImageService.GetRandomImage(category)
	if imagePath == "" {
		logger.Logger.Warn("No images available for category",
			zap.String("category", category),
			zap.String("user", i.Member.User.Username))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No %s images available", category),
			},
		})
		return
	}

	// Respond with "thinking" first
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		logger.Logger.Error("Failed to defer interaction response", zap.Error(err))
		return
	}

	// Load and send the image
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, fileName, err := h.ImageService.GetImageFile(ctx, imagePath)
	if err != nil {
		logger.Logger.Error("Failed to load image file",
			zap.String("category", category),
			zap.String("image_path", imagePath),
			zap.String("user", i.Member.User.Username),
			zap.Error(err))

		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Failed to load %s: %v", category, err),
		})
		return
	}
	defer reader.Close()

	// Send the image as a follow-up
	_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Files: []*discordgo.File{{
			Name:   fileName,
			Reader: reader,
		}},
	})

	duration := time.Since(startTime)

	if err != nil {
		logger.Logger.Error("Failed to send image",
			zap.String("category", category),
			zap.String("filename", fileName),
			zap.String("user", i.Member.User.Username),
			zap.Duration("duration", duration),
			zap.Error(err))

		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Failed to send %s: %v", category, err),
		})
	} else {
		logger.Logger.Info("Image sent successfully via slash command",
			zap.String("category", category),
			zap.String("filename", fileName),
			zap.String("user", i.Member.User.Username),
			zap.String("user_id", i.Member.User.ID),
			zap.String("channel_id", i.ChannelID),
			zap.Duration("duration", duration))
	}
}
