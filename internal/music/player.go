package music

import (
	"fmt"
	"sync"
	"time"
)

type Track struct {
	ID       string
	Title    string
	URL      string
	Duration time.Duration
	Thumbnail string
	Uploader string
}

type Player struct {
	Queue     []*Track
	Playing   bool
	Current   *Track
	mu        sync.Mutex
	volume    float64
	voiceConn *VoiceConnection
}

// VoiceConnection is a mock type representing a voice connection
// In a real implementation, this would interface with DiscordGo's voice connection
type VoiceConnection struct {
	ChannelID string
	Connected bool
}

func NewPlayer() *Player {
	return &Player{
		Queue:  make([]*Track, 0),
		volume: 1.0, // 100% volume
		voiceConn: &VoiceConnection{
			Connected: false,
		},
	}
}

func (p *Player) AddToQueue(track *Track) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.Queue = append(p.Queue, track)
}

func (p *Player) Play() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// If already playing, don't start another
	if p.Playing && p.Current != nil {
		return nil
	}
	
	if len(p.Queue) == 0 {
		return fmt.Errorf("queue is empty")
	}
	
	p.Current = p.Queue[0]
	p.Queue = p.Queue[1:]
	p.Playing = true
	
	// In a real implementation, this would:
	// 1. Connect to voice channel
	// 2. Stream the audio from URL
	// 3. Handle playback controls
	
	return nil
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.Playing = false
}

func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.Current != nil {
		p.Playing = true
	}
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.Playing = false
	p.Current = nil
	p.Queue = make([]*Track, 0)
}

func (p *Player) Skip() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if len(p.Queue) > 0 {
		p.Current = p.Queue[0]
		p.Queue = p.Queue[1:]
		return nil
	}
	
	// If queue is empty, stop playback
	p.Playing = false
	p.Current = nil
	return fmt.Errorf("queue is empty")
}

func (p *Player) SetVolume(volume float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}
	
	p.volume = volume
}

func (p *Player) GetVolume() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	return p.volume
}

func (p *Player) GetCurrentTrack() *Track {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	return p.Current
}

func (p *Player) GetQueue() []*Track {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Return a copy of the queue
	queueCopy := make([]*Track, len(p.Queue))
	copy(queueCopy, p.Queue)
	
	return queueCopy
}

func (p *Player) GetQueueLength() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	return len(p.Queue)
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	return p.Playing
}

func (p *Player) ClearQueue() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.Queue = make([]*Track, 0)
}

func (p *Player) RemoveFromQueue(index int) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if index < 0 || index >= len(p.Queue) {
		return fmt.Errorf("index out of range")
	}
	
	// Remove track at index
	p.Queue = append(p.Queue[:index], p.Queue[index+1:]...)
	
	return nil
}

func (p *Player) MoveInQueue(from, to int) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if from < 0 || from >= len(p.Queue) || to < 0 || to >= len(p.Queue) {
		return fmt.Errorf("index out of range")
	}
	
	if from == to {
		return nil
	}
	
	// Move track from 'from' to 'to'
	track := p.Queue[from]
	
	// Remove from original position
	p.Queue = append(p.Queue[:from], p.Queue[from+1:]...)
	
	// Insert at new position
	if to > from {
		to-- // Adjust for removed element
	}
	
	p.Queue = append(p.Queue[:to], append([]*Track{track}, p.Queue[to:]...)...)
	
	return nil
}

func (p *Player) ConnectToVoice(channelID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// In a real implementation, this would connect to the voice channel
	p.voiceConn.ChannelID = channelID
	p.voiceConn.Connected = true
	
	return nil
}

func (p *Player) DisconnectFromVoice() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// In a real implementation, this would disconnect from the voice channel
	p.voiceConn.Connected = false
	p.voiceConn.ChannelID = ""
	
	return nil
}

func (p *Player) IsConnectedToVoice() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	return p.voiceConn.Connected
}