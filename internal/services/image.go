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
)

type ImageService struct {
	categories map[string][]string
}

func NewImageService(baseDir string) (*ImageService, error) {
	service := &ImageService{
		categories: make(map[string][]string),
	}

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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
				return err
			}
			parts := strings.Split(relPath, string(filepath.Separator))
			if len(parts) >= 2 {
				category := parts[0]
				service.categories[category] = append(service.categories[category], path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("scan image directory: %w", err)
	}

	if len(service.categories) == 0 {
		return nil, fmt.Errorf("no image categories found in directory: %s", baseDir)
	}

	return service, nil
}

func (s *ImageService) GetRandomImage(category string) string {
	images, exists := s.categories[category]
	if !exists || len(images) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return images[rand.Intn(len(images))]
}

func (s *ImageService) GetImageFile(ctx context.Context, imagePath string) (io.ReadCloser, string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, "", fmt.Errorf("open image file: %w", err)
	}

	fileName := filepath.Base(imagePath)

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
