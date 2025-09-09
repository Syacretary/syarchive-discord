package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	apiKey   string
	engineID string
	httpClient *http.Client
}

type SearchResult struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

type ScrapedContent struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewClient(apiKey, engineID string) *Client {
	return &Client{
		apiKey:   apiKey,
		engineID: engineID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Search(query string) (*SearchResult, error) {
	baseURL := "https://www.googleapis.com/customsearch/v1"
	
	params := url.Values{}
	params.Add("key", c.apiKey)
	params.Add("cx", c.engineID)
	params.Add("q", query)
	params.Add("num", "4") // Get top 4 results
	
	searchURL := baseURL + "?" + params.Encode()
	
	resp, err := c.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search API returned status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}
	
	return &result, nil
}

func (c *Client) ScrapeContent(url string) (*ScrapedContent, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("URL returned status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read page content: %w", err)
	}
	
	// Simple content extraction (in a real implementation, you might want to use
	// a proper HTML parser like goquery)
	content := extractTextFromHTML(string(body))
	
	return &ScrapedContent{
		URL:     url,
		Content: content,
	}, nil
}

// Simple HTML text extraction (in a real implementation, use a proper HTML parser)
func extractTextFromHTML(html string) string {
	// This is a very basic implementation
	// In a real-world scenario, you would use a library like goquery
	
	// Remove script and style tags
	content := html
	
	// Simple approach to extract text content
	// This is just a placeholder - a real implementation would be more sophisticated
	
	// Limit content length to prevent token overflow
	if len(content) > 2000 {
		content = content[:2000]
	}
	
	return content
}