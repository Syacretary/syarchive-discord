package ytdlp

import (
	"testing"
)

func TestNewDownloader(t *testing.T) {
	downloader := NewDownloader()

	if downloader.maxConcurrent != 3 {
		t.Errorf("Expected maxConcurrent to be 3, got %d", downloader.maxConcurrent)
	}
}

func TestDownloadOptions(t *testing.T) {
	opts := DownloadOptions{
		URL:      "https://example.com/video.mp4",
		Format:   "mp4",
		Audio:    true,
		NoCookie: true,
	}

	if opts.URL != "https://example.com/video.mp4" {
		t.Errorf("Expected URL to be 'https://example.com/video.mp4', got '%s'", opts.URL)
	}

	if opts.Format != "mp4" {
		t.Errorf("Expected Format to be 'mp4', got '%s'", opts.Format)
	}

	if !opts.Audio {
		t.Error("Expected Audio to be true, got false")
	}

	if !opts.NoCookie {
		t.Error("Expected NoCookie to be true, got false")
	}
}

func TestVideoInfoStruct(t *testing.T) {
	info := VideoInfo{
		Title:       "Test Video",
		Duration:    300,
		Uploader:    "Test User",
		ViewCount:   1000,
		LikeCount:   50,
		Description: "This is a test video",
		Thumbnail:   "https://example.com/thumbnail.jpg",
	}

	if info.Title != "Test Video" {
		t.Errorf("Expected Title to be 'Test Video', got '%s'", info.Title)
	}

	if info.Duration != 300 {
		t.Errorf("Expected Duration to be 300, got %d", info.Duration)
	}

	if info.Uploader != "Test User" {
		t.Errorf("Expected Uploader to be 'Test User', got '%s'", info.Uploader)
	}

	if info.ViewCount != 1000 {
		t.Errorf("Expected ViewCount to be 1000, got %d", info.ViewCount)
	}

	if info.LikeCount != 50 {
		t.Errorf("Expected LikeCount to be 50, got %d", info.LikeCount)
	}

	if info.Description != "This is a test video" {
		t.Errorf("Expected Description to be 'This is a test video', got '%s'", info.Description)
	}

	if info.Thumbnail != "https://example.com/thumbnail.jpg" {
		t.Errorf("Expected Thumbnail to be 'https://example.com/thumbnail.jpg', got '%s'", info.Thumbnail)
	}
}