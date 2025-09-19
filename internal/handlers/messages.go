package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"wooper-bot/internal/services"

	"github.com/bwmarrin/discordgo"
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
	
	// Check if message starts with ! and has a valid category
	if strings.HasPrefix(content, "!") {
		category := strings.TrimPrefix(content, "!")
		if h.ImageService.HasCategory(category) {
			imagePath := h.ImageService.GetRandomImage(category)
			if imagePath == "" {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("no %s images available", category))
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			reader, fileName, err := h.ImageService.GetImageFile(ctx, imagePath)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("failed to load %s: %v", category, err))
				return
			}
			defer reader.Close()
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Files: []*discordgo.File{{
				Name:   fileName,
				Reader: reader,
			}}})
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("failed to send %s: %v", category, err))
			}
		} else if category == "help" || category == "list" {
			// Show available categories
			categories := h.ImageService.GetAvailableCategories()
			if len(categories) == 0 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "no image categories available")
				return
			}
			message := "Available image categories:\n"
			for _, cat := range categories {
				count := h.ImageService.GetImageCount(cat)
				message += fmt.Sprintf("• `!%s` (%d images)\n", cat, count)
			}
			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}
	}
}
