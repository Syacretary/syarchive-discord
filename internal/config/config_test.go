package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary .env file for testing
	testEnv := `.env.test`
	testContent := `DISCORD_TOKEN=test_token
OPENROUTER_API_KEY=test_api_key
BOT_PREFIX=!
MAX_CONCURRENT_DOWNLOADS=5
MAX_FILE_SIZE=50`

	err := os.WriteFile(testEnv, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(testEnv)

	// Set the config file to our test file
	os.Setenv("CONFIG_FILE", testEnv)
	defer os.Unsetenv("CONFIG_FILE")

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test values
	if config.DiscordToken != "test_token" {
		t.Errorf("Expected DiscordToken to be 'test_token', got '%s'", config.DiscordToken)
	}

	if config.OpenRouterAPIKey != "test_api_key" {
		t.Errorf("Expected OpenRouterAPIKey to be 'test_api_key', got '%s'", config.OpenRouterAPIKey)
	}

	if config.BotPrefix != "!" {
		t.Errorf("Expected BotPrefix to be '!', got '%s'", config.BotPrefix)
	}

	if config.MaxConcurrentDownloads != 5 {
		t.Errorf("Expected MaxConcurrentDownloads to be 5, got %d", config.MaxConcurrentDownloads)
	}

	if config.MaxFileSize != 50 {
		t.Errorf("Expected MaxFileSize to be 50, got %d", config.MaxFileSize)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Test loading config with missing values (should use defaults)
	testEnv := `.env.test_defaults`
	testContent := `DISCORD_TOKEN=test_token`

	err := os.WriteFile(testEnv, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(testEnv)

	// Set the config file to our test file
	os.Setenv("CONFIG_FILE", testEnv)
	defer os.Unsetenv("CONFIG_FILE")

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test default values
	if config.BotPrefix != "!" {
		t.Errorf("Expected BotPrefix to be '!' (default), got '%s'", config.BotPrefix)
	}

	if config.MaxConcurrentDownloads != 3 {
		t.Errorf("Expected MaxConcurrentDownloads to be 3 (default), got %d", config.MaxConcurrentDownloads)
	}

	if config.MaxFileSize != 100 {
		t.Errorf("Expected MaxFileSize to be 100 (default), got %d", config.MaxFileSize)
	}
}