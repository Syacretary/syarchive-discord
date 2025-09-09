package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"sync"
	"time"
	"encoding/json"
	"io"
	"net/http"

	"discord-bot/internal/config"
	"discord-bot/internal/openrouter"
	"discord-bot/internal/ytdlp"
	"discord-bot/internal/music"
	"discord-bot/internal/security"
	"discord-bot/internal/search"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session             *discordgo.Session
	OpenRouter          *openrouter.Client
	MusicPlayer         *music.Player
	Downloader          *ytdlp.Downloader
	Config              *config.Config
	RateLimiter         *security.RateLimiter
	MessageCounters     map[string]int // channelID -> message count
	MessageHistory      map[string][]MessageHistory // channelID -> messages
	VoiceChannelManager *VoiceChannelManager
	SearchClient        *search.Client
	mu                  sync.Mutex
	LastChannelID       string // To store the last channel ID for tool responses
}

type MessageHistory struct {
	Author    string
	Content   string
	Timestamp time.Time
}

type VoiceChannelManager struct {
	CreatedChannels map[string]string // userID -> channelID
}

func NewVoiceChannelManager() *VoiceChannelManager {
	return &VoiceChannelManager{
		CreatedChannels: make(map[string]string),
	}
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Discord session
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Initialize bot components
	bot := &Bot{
		Session:             dg,
		Config:              cfg,
		OpenRouter:          openrouter.NewClient(cfg.OpenRouterAPIKey),
		Downloader:          ytdlp.NewDownloader(),
		MusicPlayer:         music.NewPlayer(),
		RateLimiter:         security.NewRateLimiter(5, 60), // 5 requests per minute
		MessageCounters:     make(map[string]int),
		MessageHistory:      make(map[string][]MessageHistory),
		VoiceChannelManager: NewVoiceChannelManager(),
		SearchClient:        search.NewClient(cfg.GoogleSearchAPIKey, cfg.GoogleSearchEngineID),
		LastChannelID:       "",
	}

	// Register event handlers
	dg.AddHandler(bot.messageCreate)
	dg.AddHandler(bot.ready)
	dg.AddHandler(bot.voiceStateUpdate)

	// Open WebSocket connection
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Bot is ready as %v\n", event.User.Username)
}

func (b *Bot) voiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// Handle join to create voice channel
	b.handleVoiceStateUpdate(s, vs)
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Handle private messages (DMs) - forward directly to AI
	if m.GuildID == "" {
		b.handlePrivateMessage(s, m)
		return
	}

	// Handle messages in guild channels
	b.handleGuildMessage(s, m)
}

func (b *Bot) handlePrivateMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Forward private messages directly to AI without any command prefix
	question := m.Content
	
	// Send typing indicator
	s.ChannelTyping(m.ChannelID)

	// Call OpenRouter API with tools support
	messages := []openrouter.Message{
		{Role: "user", Content: question},
	}

	// Use openrouter/sonoma-dusk-alpha which supports tools
	response, err := b.OpenRouter.ChatCompletionWithTools("openrouter/sonoma-dusk-alpha", messages, openrouter.AvailableTools)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error calling AI API: %v", err))
		return
	}

	b.handleAIResponse(s, m.ChannelID, response)
}

func (b *Bot) handleGuildMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if message starts with bot prefix "/"
	if len(m.Content) >= 1 && m.Content[0] == '/' {
		// Handle command messages
		b.handleCommandMessage(s, m)
	} else {
		// Handle regular messages (for proactive AI responses)
		b.handleRegularMessage(s, m)
	}
}

func (b *Bot) handleCommandMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Apply rate limiting
	if !b.RateLimiter.IsAllowed(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "You are being rate limited. Please wait before sending more commands.")
		return
	}

	// Sanitize input
	content := security.SanitizeInput(m.Content)

	// Remove prefix and split into command and arguments
	commandParts := strings.Fields(content[1:]) // Remove the "/" prefix
	if len(commandParts) == 0 {
		return
	}

	command := strings.ToLower(commandParts[0])
	args := commandParts[1:]

	// Process commands
	switch command {
	case "ai", "ask":
		b.handleAICommand(s, m, args)
	case "download", "dl":
		b.handleDownloadCommand(s, m, args)
	case "play":
		b.handlePlayCommand(s, m, args)
	case "pause":
		b.handlePauseCommand(s, m)
	case "resume":
		b.handleResumeCommand(s, m)
	case "skip", "next":
		b.handleSkipCommand(s, m)
	case "stop":
		b.handleStopCommand(s, m)
	case "queue":
		b.handleQueueCommand(s, m)
	case "volume":
		b.handleVolumeCommand(s, m, args)
	case "help":
		b.handleHelpCommand(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown command: %s. Type /help for available commands.", command))
	}
}

func (b *Bot) handleRegularMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Add message to history
	b.MessageHistory[m.ChannelID] = append(b.MessageHistory[m.ChannelID], MessageHistory{
		Author:    m.Author.Username,
		Content:   m.Content,
		Timestamp: time.Now(),
	})

	// Keep only last 20 messages
	if len(b.MessageHistory[m.ChannelID]) > 20 {
		b.MessageHistory[m.ChannelID] = b.MessageHistory[m.ChannelID][1:]
	}

	// Increment message counter
	b.MessageCounters[m.ChannelID]++

	// Check if we should send a proactive AI response (every 10 messages)
	if b.MessageCounters[m.ChannelID] >= 10 {
		b.MessageCounters[m.ChannelID] = 0
		go b.sendProactiveAIResponse(s, m.ChannelID)
	}
}

func (b *Bot) sendProactiveAIResponse(s *discordgo.Session, channelID string) {
	b.mu.Lock()
	history := make([]MessageHistory, len(b.MessageHistory[channelID]))
	copy(history, b.MessageHistory[channelID])
	b.mu.Unlock()

	if len(history) == 0 {
		return
	}

	// Format the conversation history for the AI
	conversation := "Recent conversation in the server:\n"
	for _, msg := range history {
		conversation += fmt.Sprintf("%s: %s\n", msg.Author, msg.Content)
	}
	conversation += "\nAs an AI assistant, please provide a relevant comment or question to join the conversation naturally."

	// Send typing indicator
	s.ChannelTyping(channelID)

	// Call OpenRouter API with tools support
	messages := []openrouter.Message{
		{Role: "user", Content: conversation},
	}

	// Use openrouter/sonoma-dusk-alpha which supports tools
	response, err := b.OpenRouter.ChatCompletionWithTools("openrouter/sonoma-dusk-alpha", messages, openrouter.AvailableTools)
	if err != nil {
		// Don't send error message for proactive responses
		return
	}

	b.handleAIResponse(s, channelID, response)
}

func (b *Bot) handleAICommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a question for the AI.")
		return
	}

	question := strings.Join(args, " ")
	
	// Send typing indicator
	s.ChannelTyping(m.ChannelID)

	// Call OpenRouter API with tools support
	messages := []openrouter.Message{
		{Role: "user", Content: question},
	}

	// Use openrouter/sonoma-dusk-alpha which supports tools
	response, err := b.OpenRouter.ChatCompletionWithTools("openrouter/sonoma-dusk-alpha", messages, openrouter.AvailableTools)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error calling AI API: %v", err))
		return
	}

	b.handleAIResponse(s, m.ChannelID, response)
}

func (b *Bot) handleAIResponse(s *discordgo.Session, channelID string, response *openrouter.ChatResponse) {
	// Store the channel ID for tool responses
	b.LastChannelID = channelID

	if len(response.Choices) > 0 {
		choice := response.Choices[0]
		
		// Check if the AI wants to call any tools
		if len(choice.ToolCalls) > 0 {
			// Process tool calls
			for _, toolCall := range choice.ToolCalls {
				result := b.executeTool(toolCall.Function.Name, toolCall.Function.Arguments)
				
				// Send the result back to the AI
				messages := []openrouter.Message{
					{Role: "user", Content: choice.Message.Content},
					{Role: "assistant", Content: "", ToolCalls: choice.ToolCalls},
					{Role: "tool", Content: result, Name: toolCall.Function.Name},
				}
				
				// Get final response from AI
				finalResponse, err := b.OpenRouter.ChatCompletionWithTools("openrouter/sonoma-dusk-alpha", messages, openrouter.AvailableTools)
				if err != nil {
					s.ChannelMessageSend(channelID, fmt.Sprintf("Error getting final response: %v", err))
					return
				}
				
				if len(finalResponse.Choices) > 0 && finalResponse.Choices[0].Message.Content != "" {
					s.ChannelMessageSend(channelID, finalResponse.Choices[0].Message.Content)
				}
			}
		} else if choice.Message.Content != "" {
			// Regular response without tool calls
			s.ChannelMessageSend(channelID, choice.Message.Content)
		} else {
			s.ChannelMessageSend(channelID, "No response from AI.")
		}
	} else {
		s.ChannelMessageSend(channelID, "No response from AI.")
	}
}

func (b *Bot) executeTool(name, arguments string) string {
	switch name {
	case "download_video":
		return b.executeDownloadVideo(arguments)
	case "play_music":
		return b.executePlayMusic(arguments)
	case "get_video_info":
		return b.executeGetVideoInfo(arguments)
	case "search_web":
		return b.executeSearchWeb(arguments)
	default:
		return fmt.Sprintf("Unknown tool: %s", name)
	}
}

func (b *Bot) executeDownloadVideo(arguments string) string {
	var args struct {
		URL    string `json:"url"`
		Format string `json:"format"`
	}
	
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	
	audioOnly := args.Format == "audio"
	
	opts := ytdlp.DownloadOptions{
		URL:      args.URL,
		Audio:    audioOnly,
		NoCookie: true,
	}
	
	filename, err := b.Downloader.DownloadVideo(opts)
	if err != nil {
		return fmt.Sprintf("Error downloading: %v", err)
	}
	
	return fmt.Sprintf("Download completed: %s", filename)
}

func (b *Bot) executePlayMusic(arguments string) string {
	var args struct {
		URL string `json:"url"`
	}
	
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	
	track := &music.Track{
		Title: args.URL,
		URL:   args.URL,
	}
	
	b.MusicPlayer.AddToQueue(track)
	return fmt.Sprintf("Added to queue: %s", args.URL)
}

func (b *Bot) executeGetVideoInfo(arguments string) string {
	var args struct {
		URL string `json:"url"`
	}
	
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	
	info, err := b.Downloader.GetInfo(args.URL)
	if err != nil {
		return fmt.Sprintf("Error getting video info: %v", err)
	}
	
	return fmt.Sprintf("Video info: %+v", info)
}

func (b *Bot) executeSearchWeb(arguments string) string {
	var args struct {
		Query string `json:"query"`
	}
	
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	
	// Step 1: Perform Google Custom Search
	searchResults, err := b.SearchClient.Search(args.Query)
	if err != nil {
		return fmt.Sprintf("Error performing search: %v", err)
	}
	
	if len(searchResults.Items) == 0 {
		return "No search results found."
	}
	
	// Step 2: Scrape content from top 4 websites
	scrapedContents := make([]string, 0, len(searchResults.Items))
	for _, item := range searchResults.Items {
		content, err := b.SearchClient.ScrapeContent(item.Link)
		if err != nil {
			// If scraping fails, use the snippet instead
			scrapedContents = append(scrapedContents, fmt.Sprintf("Title: %s\nURL: %s\nContent: %s", item.Title, item.Link, item.Snippet))
		} else {
			scrapedContents = append(scrapedContents, fmt.Sprintf("Title: %s\nURL: %s\nContent: %s", item.Title, item.Link, content.Content))
		}
	}
	
	// Step 3: Combine all scraped content
	allContent := strings.Join(scrapedContents, "\n\n---\n\n")
	
	// Step 4: Ask AI to summarize the search results
	summary, err := b.summarizeSearchResults(allContent)
	if err != nil {
		return fmt.Sprintf("Error summarizing results: %v", err)
	}
	
	return fmt.Sprintf("**Search Results for '%s':**\n\n%s", args.Query, summary)
}

func (b *Bot) summarizeSearchResults(content string) (string, error) {
	// Send typing indicator
	if b.LastChannelID != "" {
		b.Session.ChannelTyping(b.LastChannelID)
	}
	
	// Ask AI to summarize the search results
	messages := []openrouter.Message{
		{Role: "user", Content: fmt.Sprintf("Tolong rangkum hasil web search ini dengan rapi:\n\n%s", content)},
	}
	
	// Use openrouter/sonoma-dusk-alpha for summarization
	response, err := b.OpenRouter.ChatCompletion("openrouter/sonoma-dusk-alpha", messages)
	if err != nil {
		return "", fmt.Errorf("error calling AI API for summarization: %w", err)
	}
	
	if len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		return response.Choices[0].Message.Content, nil
	}
	
	return "Unable to generate summary.", nil
}

func (b *Bot) handleDownloadCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a URL to download.")
		return
	}

	url := args[0]
	if !security.ValidateURL(url) {
		s.ChannelMessageSend(m.ChannelID, "Invalid URL provided.")
		return
	}

	// Send processing message
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Downloading from %s...", url))

	// Determine if audio only
	audioOnly := false
	for _, arg := range args {
		if arg == "-a" || arg == "--audio" {
			audioOnly = true
			break
		}
	}

	opts := ytdlp.DownloadOptions{
		URL:      url,
		Audio:    audioOnly,
		NoCookie: true,
	}

	filename, err := b.Downloader.DownloadVideo(opts)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error downloading: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Download completed: %s", filename))
}

func (b *Bot) handlePlayCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a URL to play.")
		return
	}

	url := args[0]
	if !security.ValidateURL(url) {
		s.ChannelMessageSend(m.ChannelID, "Invalid URL provided.")
		return
	}

	// In a real implementation, we would:
	// 1. Connect to voice channel
	// 2. Add track to queue
	// 3. Start playing if not already playing
	
	track := &music.Track{
		Title: url,
		URL:   url,
	}
	
	b.MusicPlayer.AddToQueue(track)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added to queue: %s", url))
}

func (b *Bot) handlePauseCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.MusicPlayer.Pause()
	s.ChannelMessageSend(m.ChannelID, "Playback paused.")
}

func (b *Bot) handleResumeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.MusicPlayer.Play()
	s.ChannelMessageSend(m.ChannelID, "Playback resumed.")
}

func (b *Bot) handleSkipCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.MusicPlayer.Skip()
	s.ChannelMessageSend(m.ChannelID, "Skipped to next track.")
}

func (b *Bot) handleStopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.MusicPlayer.Stop()
	s.ChannelMessageSend(m.ChannelID, "Playback stopped and queue cleared.")
}

func (b *Bot) handleQueueCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	queue := b.MusicPlayer.GetQueue()
	if len(queue) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Queue is empty.")
		return
	}

	message := "Current queue:\n"
	for i, track := range queue {
		message += fmt.Sprintf("%d. %s\n", i+1, track.Title)
	}
	
	s.ChannelMessageSend(m.ChannelID, message)
}

func (b *Bot) handleVolumeCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		volume := b.MusicPlayer.GetVolume() * 100
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Current volume: %.0f%%", volume))
		return
	}

	// Parse volume level
	var volume int
	_, err := fmt.Sscanf(args[0], "%d", &volume)
	if err != nil || volume < 0 || volume > 100 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a valid volume level between 0 and 100.")
		return
	}

	b.MusicPlayer.SetVolume(float64(volume) / 100.0)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Volume set to %d%%", volume))
}

func (b *Bot) handleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpText := fmt.Sprintf("Available commands:\n"+
		"/help - Show this help message\n"+
		"/ai <question> - Ask the AI a question\n"+
		"/download <url> [-a] - Download video/audio from URL (-a for audio only)\n"+
		"/play <url> - Play audio from URL\n"+
		"/pause - Pause playback\n"+
		"/resume - Resume playback\n"+
		"/skip - Skip to next track\n"+
		"/stop - Stop playback and clear queue\n"+
		"/queue - Show current queue\n"+
		"/volume [level] - Show or set volume (0-100)")

	s.ChannelMessageSend(m.ChannelID, helpText)
}

// Voice channel management functions
func (b *Bot) handleVoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// Check if this is a join event to a specific "Join to Create" channel
	// You would need to configure this channel ID in your config
	if vs.ChannelID != "" && b.isJoinToCreateChannel(vs.ChannelID) {
		b.createVoiceChannelForUser(s, vs)
	}
	
	// Check if a user left a created channel and it's now empty
	if vs.BeforeUpdate != nil && vs.BeforeUpdate.ChannelID != "" {
		b.checkAndDeleteEmptyChannel(s, vs.BeforeUpdate.ChannelID)
	}
}

func (b *Bot) isJoinToCreateChannel(channelID string) bool {
	// In a real implementation, you would check this against a configured channel ID
	// For now, we'll return false to avoid accidental channel creation
	return false
}

func (b *Bot) createVoiceChannelForUser(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// Get the guild ID
	guildID := vs.GuildID
	
	// Get user information
	user, err := s.User(vs.UserID)
	if err != nil {
		return
	}
	
	// Create a new voice channel
	channel, err := s.GuildChannelCreate(guildID, user.Username+"'s Channel", discordgo.ChannelTypeGuildVoice)
	if err != nil {
		return
	}
	
	// Move the user to the new channel
	err = s.GuildMemberMove(guildID, vs.UserID, &channel.ID)
	if err != nil {
		// If we can't move the user, delete the channel we just created
		s.ChannelDelete(channel.ID)
		return
	}
	
	// Store the mapping
	b.mu.Lock()
	b.VoiceChannelManager.CreatedChannels[vs.UserID] = channel.ID
	b.mu.Unlock()
}

func (b *Bot) checkAndDeleteEmptyChannel(s *discordgo.Session, channelID string) {
	// Check if this channel was created by our bot
	b.mu.Lock()
	userID := ""
	for uID, chID := range b.VoiceChannelManager.CreatedChannels {
		if chID == channelID {
			userID = uID
			break
		}
	}
	b.mu.Unlock()
	
	if userID == "" {
		// This channel wasn't created by us
		return
	}
	
	// Get the channel to check if it's empty
	channel, err := s.Channel(channelID)
	if err != nil {
		return
	}
	
	// Check if the channel is empty
	if len(channel.VoiceStates) == 0 {
		// Delete the channel
		s.ChannelDelete(channelID)
		
		// Remove from our mapping
		b.mu.Lock()
		delete(b.VoiceChannelManager.CreatedChannels, userID)
		b.mu.Unlock()
	}
}