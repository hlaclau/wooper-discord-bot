package bot

import (
	"context"
	"testing"
	"time"
)

// TestNew tests the New function for creating Discord bot instances.
func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "valid token",
			token:       "valid-token-123",
			expectError: false,
		},
		{
			name:        "empty token",
			token:       "",
			expectError: false, // DiscordGo doesn't validate empty tokens at creation
		},
		{
			name:        "invalid token format",
			token:       "invalid format with spaces",
			expectError: false, // DiscordGo doesn't validate token format at creation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot, err := New(tt.token)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if bot != nil {
					t.Errorf("Expected nil bot but got %v", bot)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if bot == nil {
					t.Errorf("Expected bot but got nil")
				}
				if bot.session == nil {
					t.Errorf("Expected session but got nil")
				}
			}
		})
	}
}

// TestBot_AddHandler tests the AddHandler method of the Bot struct.
func TestBot_AddHandler(t *testing.T) {
	bot, err := New("test-token")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Test that AddHandler returns a function
	removeHandler := bot.AddHandler(func() {})
	if removeHandler == nil {
		t.Errorf("Expected remove handler function but got nil")
	}

	// Test that the returned function can be called
	removeHandler()
}

func TestBot_Start(t *testing.T) {
	// Note: This test is limited because we can't easily test Discord session
	// without mocking or actual Discord connection
	bot, err := New("test-token")
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	// Create a context that will be cancelled quickly
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// This will likely fail because we don't have a real Discord token,
	// but we can test that the function doesn't panic
	err = bot.Start(ctx)

	// We expect an error because we can't actually connect to Discord
	if err == nil {
		t.Log("Note: Start() succeeded unexpectedly (this might be due to test environment)")
	}
}
