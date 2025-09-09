package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message   Message    `json:"message"`
		ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	} `json:"choices"`
}

// Supported models
const (
	ModelClaude3Sonnet  = "claude-3-sonnet"
	ModelLLaVA          = "llava-v1.5-7b-4bit"
	ModelMistral7B      = "mistralai/mistral-7b-instruct"
	ModelZephyr7B       = "huggingfaceh4/zephyr-7b-beta"
	ModelMythoMax       = "gryphe/mythomax-l2-13b"
	ModelToppyM         = "undi95/toppy-m-7b"
	ModelOpenChat       = "openchat/openchat-7b"
	ModelSonomaDusk     = "openrouter/sonoma-dusk-alpha" // Cepat, akurat, dan kuat untuk merespon pengguna di Discord
)

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) ChatCompletion(model string, messages []Message) (*ChatResponse, error) {
	request := ChatRequest{
		Model:    model,
		Messages: messages,
	}

	return c.sendRequest(request)
}

func (c *Client) ChatCompletionWithTools(model string, messages []Message, tools []Tool) (*ChatResponse, error) {
	request := ChatRequest{
		Model:    model,
		Messages: messages,
		Tools:    tools,
	}

	return c.sendRequest(request)
}

func (c *Client) sendRequest(request ChatRequest) (*ChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chatResp, nil
}

func (c *Client) SupportedModels() []string {
	return []string{
		ModelClaude3Sonnet,
		ModelLLaVA,
		ModelMistral7B,
		ModelZephyr7B,
		ModelMythoMax,
		ModelToppyM,
		ModelOpenChat,
		ModelSonomaDusk,
	}
}

// Tools that can be used by the AI
var AvailableTools = []Tool{
	{
		Type: "function",
		Function: Function{
			Name:        "download_video",
			Description: "Download a video or audio from a given URL",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL of the video to download",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"description": "The format to download (video or audio)",
						"enum":        []string{"video", "audio"},
					},
				},
				"required": []string{"url"},
			},
		},
	},
	{
		Type: "function",
		Function: Function{
			Name:        "play_music",
			Description: "Play music from a given URL",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL of the music to play",
					},
				},
				"required": []string{"url"},
			},
		},
	},
	{
		Type: "function",
		Function: Function{
			Name:        "get_video_info",
			Description: "Get information about a video from a given URL",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL of the video to get information about",
					},
				},
				"required": []string{"url"},
			},
		},
	},
	{
		Type: "function",
		Function: Function{
			Name:        "search_web",
			Description: "Search the web for information",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
				},
				"required": []string{"query"},
			},
		},
	},
}