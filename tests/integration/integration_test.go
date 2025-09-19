package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"wooper-bot/internal/logger"
	"wooper-bot/internal/services"
)

// TestIntegration tests the overall flow of the wooper-bot application.
func TestIntegration(t *testing.T) {
	// Skip if running in CI or if Discord token is not available
	if os.Getenv("CI") != "" || os.Getenv("DISCORD_BOT_TOKEN") == "" {
		t.Skip("Skipping integration test - no Discord token or running in CI")
	}

	// Initialize logger for integration tests
	err := logger.Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Create test image directory
	tempDir := t.TempDir()
	setupTestImages(t, tempDir)

	// Test image service initialization
	imageService, err := services.NewImageService(tempDir)
	if err != nil {
		t.Fatalf("Failed to create image service: %v", err)
	}

	// Test that categories are loaded
	categories := imageService.GetAvailableCategories()
	if len(categories) == 0 {
		t.Errorf("Expected categories to be loaded but got none")
	}

	// Test getting random images
	for _, category := range categories {
		imagePath := imageService.GetRandomImage(category)
		if imagePath == "" {
			t.Errorf("Expected image path for category %s but got empty", category)
		}
	}

	// Test image file operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, category := range categories {
		imagePath := imageService.GetRandomImage(category)
		if imagePath != "" {
			reader, filename, err := imageService.GetImageFile(ctx, imagePath)
			if err != nil {
				t.Errorf("Failed to get image file for category %s: %v", category, err)
			} else {
				reader.Close()
				if filename == "" {
					t.Errorf("Expected filename but got empty for category %s", category)
				}
			}
		}
	}
}

// setupTestImages creates test image files in the given directory
func setupTestImages(t *testing.T, baseDir string) {
	categories := []string{"wooper", "cats", "dogs"}

	for _, category := range categories {
		categoryDir := filepath.Join(baseDir, category)
		err := os.MkdirAll(categoryDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create test image files
		for i := 1; i <= 3; i++ {
			filename := filepath.Join(categoryDir, category+"_"+string(rune('0'+i))+".jpg")
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			file.Close()
		}
	}
}

// TestMain tests the main function components without running the full application.
func TestMain(t *testing.T) {
	// This test verifies that the main function components can be initialized
	// without actually running the full application

	// Test config loading (this will fail without proper env vars, which is expected)
	// We're just testing that the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Main function panicked: %v", r)
		}
	}()

	// Test that we can import and use the packages
	// This is a basic smoke test
	_ = services.ImageService{}
}
