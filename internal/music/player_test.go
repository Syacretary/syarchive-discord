package music

import (
	"testing"
	"time"
)

func TestNewPlayer(t *testing.T) {
	player := NewPlayer()

	if player.volume != 1.0 {
		t.Errorf("Expected volume to be 1.0, got %f", player.volume)
	}

	if len(player.Queue) != 0 {
		t.Errorf("Expected queue to be empty, got length %d", len(player.Queue))
	}

	if player.Playing {
		t.Error("Expected Playing to be false, got true")
	}

	if player.Current != nil {
		t.Error("Expected Current to be nil, got non-nil")
	}

	if player.voiceConn == nil {
		t.Error("Expected voiceConn to be initialized, got nil")
	}

	if player.voiceConn.Connected {
		t.Error("Expected voiceConn.Connected to be false, got true")
	}
}

func TestTrackStruct(t *testing.T) {
	track := Track{
		ID:        "123",
		Title:     "Test Song",
		URL:       "https://example.com/song.mp3",
		Duration:  180 * time.Second,
		Thumbnail: "https://example.com/thumbnail.jpg",
		Uploader:  "Test Artist",
	}

	if track.ID != "123" {
		t.Errorf("Expected ID to be '123', got '%s'", track.ID)
	}

	if track.Title != "Test Song" {
		t.Errorf("Expected Title to be 'Test Song', got '%s'", track.Title)
	}

	if track.URL != "https://example.com/song.mp3" {
		t.Errorf("Expected URL to be 'https://example.com/song.mp3', got '%s'", track.URL)
	}

	if track.Duration != 180*time.Second {
		t.Errorf("Expected Duration to be 180s, got %v", track.Duration)
	}

	if track.Thumbnail != "https://example.com/thumbnail.jpg" {
		t.Errorf("Expected Thumbnail to be 'https://example.com/thumbnail.jpg', got '%s'", track.Thumbnail)
	}

	if track.Uploader != "Test Artist" {
		t.Errorf("Expected Uploader to be 'Test Artist', got '%s'", track.Uploader)
	}
}

func TestSetVolume(t *testing.T) {
	player := NewPlayer()

	// Test setting volume within range
	player.SetVolume(0.5)
	if player.GetVolume() != 0.5 {
		t.Errorf("Expected volume to be 0.5, got %f", player.GetVolume())
	}

	// Test setting volume above range
	player.SetVolume(1.5)
	if player.GetVolume() != 1.0 {
		t.Errorf("Expected volume to be capped at 1.0, got %f", player.GetVolume())
	}

	// Test setting volume below range
	player.SetVolume(-0.5)
	if player.GetVolume() != 0.0 {
		t.Errorf("Expected volume to be capped at 0.0, got %f", player.GetVolume())
	}
}

func TestQueueOperations(t *testing.T) {
	player := NewPlayer()

	// Test adding to queue
	track1 := &Track{Title: "Song 1", URL: "url1"}
	track2 := &Track{Title: "Song 2", URL: "url2"}

	player.AddToQueue(track1)
	player.AddToQueue(track2)

	queue := player.GetQueue()
	if len(queue) != 2 {
		t.Errorf("Expected queue length to be 2, got %d", len(queue))
	}

	if queue[0].Title != "Song 1" {
		t.Errorf("Expected first track title to be 'Song 1', got '%s'", queue[0].Title)
	}

	if queue[1].Title != "Song 2" {
		t.Errorf("Expected second track title to be 'Song 2', got '%s'", queue[1].Title)
	}
}