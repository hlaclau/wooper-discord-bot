package services

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wooper-bot/internal/logger"

	"go.uber.org/zap"
)

type ImageService struct {
	categories map[string][]string
}

func NewImageService(baseDir string) (*ImageService, error) {
	logger.Logger.Info("Initializing image service", zap.String("base_dir", baseDir))

	service := &ImageService{
		categories: make(map[string][]string),
	}

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Logger.Error("Error walking directory", zap.String("path", path), zap.Error(err))
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".webp" {
			// Extract category from path (e.g., img/wooper/image.jpg -> wooper)
			relPath, err := filepath.Rel(baseDir, path)
			if err != nil {
				logger.Logger.Error("Error getting relative path", zap.String("path", path), zap.Error(err))
				return err
			}
			parts := strings.Split(relPath, string(filepath.Separator))
			if len(parts) >= 2 {
				category := parts[0]
				service.categories[category] = append(service.categories[category], path)
				logger.Logger.Debug("Found image",
					zap.String("category", category),
					zap.String("file", filepath.Base(path)),
					zap.String("path", path))
			}
		}
		return nil
	})

	if err != nil {
		logger.Logger.Error("Error scanning image directory", zap.Error(err))
		return nil, fmt.Errorf("scan image directory: %w", err)
	}

	if len(service.categories) == 0 {
		logger.Logger.Error("No image categories found", zap.String("base_dir", baseDir))
		return nil, fmt.Errorf("no image categories found in directory: %s", baseDir)
	}

	// Log summary of loaded categories
	for category, images := range service.categories {
		logger.Logger.Info("Loaded image category",
			zap.String("category", category),
			zap.Int("count", len(images)))
	}

	logger.Logger.Info("Image service initialized successfully",
		zap.Int("total_categories", len(service.categories)))

	return service, nil
}

func (s *ImageService) GetRandomImage(category string) string {
	images, exists := s.categories[category]
	if !exists || len(images) == 0 {
		logger.Logger.Warn("No images found for category", zap.String("category", category))
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	selectedImage := images[rand.Intn(len(images))]
	logger.Logger.Debug("Selected random image",
		zap.String("category", category),
		zap.String("image", filepath.Base(selectedImage)),
		zap.Int("total_available", len(images)))
	return selectedImage
}

func (s *ImageService) GetImageFile(ctx context.Context, imagePath string) (io.ReadCloser, string, error) {
	logger.Logger.Debug("Opening image file", zap.String("path", imagePath))

	file, err := os.Open(imagePath)
	if err != nil {
		logger.Logger.Error("Failed to open image file", zap.String("path", imagePath), zap.Error(err))
		return nil, "", fmt.Errorf("open image file: %w", err)
	}

	fileName := filepath.Base(imagePath)
	logger.Logger.Debug("Successfully opened image file",
		zap.String("filename", fileName),
		zap.String("path", imagePath))

	return file, fileName, nil
}

func (s *ImageService) GetImageCount(category string) int {
	if images, exists := s.categories[category]; exists {
		return len(images)
	}
	return 0
}

func (s *ImageService) GetAvailableCategories() []string {
	var categories []string
	for category := range s.categories {
		categories = append(categories, category)
	}
	return categories
}

func (s *ImageService) HasCategory(category string) bool {
	_, exists := s.categories[category]
	return exists
}
