package handlers

import (
	"os"
	"path/filepath"
	"testing"

	"wooper-bot/internal/logger"
	"wooper-bot/internal/services"
)

// setupTestHandler creates a message handler with test image service
func setupTestHandler(t *testing.T) *MessageHandler {
	// Initialize logger for tests
	err := logger.Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	t.Cleanup(func() {
		logger.Close()
	})

	// Create temporary test images
	tempDir := t.TempDir()

	// Create test image directories and files
	categories := []string{"wooper", "cats"}

	for _, category := range categories {
		categoryDir := filepath.Join(tempDir, category)
		err := os.MkdirAll(categoryDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create test image files
		for i := 1; i <= 2; i++ {
			filename := filepath.Join(categoryDir, category+"_"+string(rune('0'+i))+".jpg")
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			file.Close()
		}
	}

	// Create image service
	imageService, err := services.NewImageService(tempDir)
	if err != nil {
		t.Fatalf("Failed to create image service: %v", err)
	}

	// Create message handler
	handler := NewMessageHandler(imageService)

	return handler
}

// TestNewMessageHandler tests the NewMessageHandler constructor function.
func TestNewMessageHandler(t *testing.T) {
	// Create a mock image service
	imageService := &services.ImageService{}

	handler := NewMessageHandler(imageService)

	if handler == nil {
		t.Errorf("Expected handler but got nil")
	}

	if handler.ImageService != imageService {
		t.Errorf("Expected image service %v but got %v", imageService, handler.ImageService)
	}
}

// TestMessageHandler_ImageServiceIntegration tests the integration between MessageHandler and ImageService.
func TestMessageHandler_ImageServiceIntegration(t *testing.T) {
	handler := setupTestHandler(t)

	// Test that the handler has access to the image service
	if handler.ImageService == nil {
		t.Errorf("Expected image service but got nil")
	}

	// Test that the image service has categories
	categories := handler.ImageService.GetAvailableCategories()
	if len(categories) == 0 {
		t.Errorf("Expected categories but got none")
	}

	// Test that we can get random images
	for _, category := range categories {
		imagePath := handler.ImageService.GetRandomImage(category)
		if imagePath == "" {
			t.Errorf("Expected image path for category %s but got empty", category)
		}
	}
}
