package logger

import (
	"os"
	"testing"
)

// TestInit tests the logger initialization with various log levels.
func TestInit(t *testing.T) {
	// Test with different log levels
	levels := []string{"debug", "info", "warn", "error", "invalid"}

	for _, level := range levels {
		t.Run("level_"+level, func(t *testing.T) {
			// Set environment variable
			os.Setenv("LOG_LEVEL", level)
			defer os.Unsetenv("LOG_LEVEL")

			// Initialize logger
			err := Init()

			if err != nil {
				t.Errorf("Unexpected error initializing logger: %v", err)
			}

			if Logger == nil {
				t.Errorf("Expected logger to be initialized but got nil")
			}

			// Clean up
			Close()
		})
	}
}

func TestInit_DefaultLevel(t *testing.T) {
	// Unset LOG_LEVEL to test default
	os.Unsetenv("LOG_LEVEL")

	err := Init()
	if err != nil {
		t.Errorf("Unexpected error initializing logger: %v", err)
	}

	if Logger == nil {
		t.Errorf("Expected logger to be initialized but got nil")
	}

	// Clean up
	Close()
}

func TestClose(t *testing.T) {
	// Initialize logger
	err := Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Close should not panic
	Close()

	// Close again should not panic
	Close()
}

// TestLogger_Logging tests the logger's ability to output different log levels.
func TestLogger_Logging(t *testing.T) {
	// Initialize logger
	err := Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Close()

	// Test that we can log without panicking
	Logger.Info("Test info message")
	Logger.Debug("Test debug message")
	Logger.Warn("Test warning message")
	Logger.Error("Test error message")
}
