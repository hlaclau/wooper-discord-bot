package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"wooper-bot/internal/logger"
)

// setupTestLogger initializes the logger for tests
func setupTestLogger(t *testing.T) {
	err := logger.Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	t.Cleanup(func() {
		logger.Close()
	})
}

// setupTestImages creates a temporary directory structure with test images
func setupTestImages(t *testing.T) string {
	setupTestLogger(t)

	tempDir := t.TempDir()

	// Create test image directories and files
	categories := []string{"wooper", "cats", "dogs"}

	for _, category := range categories {
		categoryDir := filepath.Join(tempDir, category)
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

	return tempDir
}

// TestNewImageService tests the NewImageService function with various directory scenarios.
func TestNewImageService(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		expectError bool
	}{
		{
			name:        "valid directory with images",
			baseDir:     setupTestImages(t),
			expectError: false,
		},
		{
			name:        "non-existent directory",
			baseDir:     "/non/existent/path",
			expectError: true,
		},
		{
			name: "empty directory",
			baseDir: func() string {
				dir := t.TempDir()
				return dir
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewImageService(tt.baseDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if service != nil {
					t.Errorf("Expected nil service but got %v", service)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if service == nil {
					t.Errorf("Expected service but got nil")
				}
			}
		})
	}
}

// TestImageService_GetRandomImage tests the GetRandomImage method of ImageService.
func TestImageService_GetRandomImage(t *testing.T) {
	testDir := setupTestImages(t)
	service, err := NewImageService(testDir)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	tests := []struct {
		name           string
		category       string
		expectedResult bool
	}{
		{
			name:           "existing category",
			category:       "wooper",
			expectedResult: true,
		},
		{
			name:           "non-existing category",
			category:       "nonexistent",
			expectedResult: false,
		},
		{
			name:           "empty category",
			category:       "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetRandomImage(tt.category)

			if tt.expectedResult {
				if result == "" {
					t.Errorf("Expected non-empty result but got empty string")
				}
				// Verify the result is a valid path
				if _, err := os.Stat(result); os.IsNotExist(err) {
					t.Errorf("Result path does not exist: %s", result)
				}
			} else {
				if result != "" {
					t.Errorf("Expected empty result but got: %s", result)
				}
			}
		})
	}
}

// TestImageService_GetImageFile tests the GetImageFile method of ImageService.
func TestImageService_GetImageFile(t *testing.T) {
	testDir := setupTestImages(t)
	service, err := NewImageService(testDir)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Get a valid image path
	validImage := service.GetRandomImage("wooper")
	if validImage == "" {
		t.Fatalf("No valid image found for testing")
	}

	tests := []struct {
		name        string
		imagePath   string
		expectError bool
	}{
		{
			name:        "valid image path",
			imagePath:   validImage,
			expectError: false,
		},
		{
			name:        "non-existent file",
			imagePath:   "/non/existent/file.jpg",
			expectError: true,
		},
		{
			name:        "empty path",
			imagePath:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			reader, filename, err := service.GetImageFile(ctx, tt.imagePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if reader != nil {
					t.Errorf("Expected nil reader but got %v", reader)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if reader == nil {
					t.Errorf("Expected reader but got nil")
				} else {
					reader.Close() // Clean up
				}
				if filename == "" {
					t.Errorf("Expected filename but got empty string")
				}
			}
		})
	}
}

func TestImageService_GetImageCount(t *testing.T) {
	testDir := setupTestImages(t)
	service, err := NewImageService(testDir)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	tests := []struct {
		name     string
		category string
		expected int
	}{
		{
			name:     "existing category",
			category: "wooper",
			expected: 3,
		},
		{
			name:     "non-existing category",
			category: "nonexistent",
			expected: 0,
		},
		{
			name:     "empty category",
			category: "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetImageCount(tt.category)
			if result != tt.expected {
				t.Errorf("Expected count %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestImageService_GetAvailableCategories(t *testing.T) {
	testDir := setupTestImages(t)
	service, err := NewImageService(testDir)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	categories := service.GetAvailableCategories()
	expectedCount := 3 // wooper, cats, dogs

	if len(categories) != expectedCount {
		t.Errorf("Expected %d categories, got %d", expectedCount, len(categories))
	}

	// Check that all expected categories are present
	expectedCategories := []string{"wooper", "cats", "dogs"}
	for _, expected := range expectedCategories {
		found := false
		for _, actual := range categories {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected category %s not found in %v", expected, categories)
		}
	}
}

func TestImageService_HasCategory(t *testing.T) {
	testDir := setupTestImages(t)
	service, err := NewImageService(testDir)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	tests := []struct {
		name     string
		category string
		expected bool
	}{
		{
			name:     "existing category",
			category: "wooper",
			expected: true,
		},
		{
			name:     "non-existing category",
			category: "nonexistent",
			expected: false,
		},
		{
			name:     "empty category",
			category: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.HasCategory(tt.category)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestImageService_FileExtensions tests the ImageService's ability to filter image file extensions.
func TestImageService_FileExtensions(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test directory with various file extensions
	testDir := filepath.Join(tempDir, "test")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create files with different extensions
	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".txt", ".md"}
	for _, ext := range extensions {
		filename := filepath.Join(testDir, "test"+ext)
		file, err := os.Create(filename)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()
	}

	service, err := NewImageService(tempDir)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Should only find image files (5 out of 7 files)
	count := service.GetImageCount("test")
	expectedCount := 5 // png, jpg, jpeg, gif, webp

	if count != expectedCount {
		t.Errorf("Expected %d image files, got %d", expectedCount, count)
	}
}
