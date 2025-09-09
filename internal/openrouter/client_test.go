package openrouter

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key"
	client := NewClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey to be '%s', got '%s'", apiKey, client.apiKey)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized, got nil")
	}

	if client.httpClient.Timeout == 0 {
		t.Error("Expected httpClient timeout to be set, got 0")
	}
}

func TestSupportedModels(t *testing.T) {
	client := NewClient("test_api_key")
	models := client.SupportedModels()

	expectedModels := []string{
		ModelClaude3Sonnet,
		ModelLLaVA,
		ModelMistral7B,
		ModelZephyr7B,
	}

	if len(models) != len(expectedModels) {
		t.Errorf("Expected %d models, got %d", len(expectedModels), len(models))
	}

	// Check if all expected models are present
	for _, expected := range expectedModels {
		found := false
		for _, model := range models {
			if model == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected model '%s' not found in supported models", expected)
		}
	}
}

func TestMessageStruct(t *testing.T) {
	message := Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	if message.Role != "user" {
		t.Errorf("Expected Role to be 'user', got '%s'", message.Role)
	}

	if message.Content != "Hello, world!" {
		t.Errorf("Expected Content to be 'Hello, world!', got '%s'", message.Content)
	}
}