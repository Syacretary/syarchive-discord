package security

import (
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Permission string

const (
	PermissionAdmin    Permission = "admin"
	PermissionMusic    Permission = "music"
	PermissionDownload Permission = "download"
	PermissionAI       Permission = "ai"
)

type RateLimiter struct {
	limits map[string][]time.Time
	maxReq int
	window time.Duration
}

type UserPermissions struct {
	UserID      string
	Permissions []Permission
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limits: make(map[string][]time.Time),
		maxReq: maxRequests,
		window: window,
	}
}

func (rl *RateLimiter) IsAllowed(userID string) bool {
	now := time.Now()
	
	// Clean up old entries
	var validRequests []time.Time
	for _, reqTime := range rl.limits[userID] {
		if now.Sub(reqTime) < rl.window {
			validRequests = append(validRequests, reqTime)
		}
	}
	rl.limits[userID] = validRequests
	
	// Check if user is within limit
	if len(rl.limits[userID]) >= rl.maxReq {
		return false
	}
	
	// Add new request
	rl.limits[userID] = append(rl.limits[userID], now)
	return true
}

func (rl *RateLimiter) GetRemainingRequests(userID string) int {
	now := time.Now()
	
	// Clean up old entries
	var validRequests []time.Time
	for _, reqTime := range rl.limits[userID] {
		if now.Sub(reqTime) < rl.window {
			validRequests = append(validRequests, reqTime)
		}
	}
	rl.limits[userID] = validRequests
	
	return rl.maxReq - len(rl.limits[userID])
}

func (rl *RateLimiter) GetResetTime(userID string) time.Time {
	if len(rl.limits[userID]) == 0 {
		return time.Now()
	}
	
	// Return the time when the oldest request will expire
	oldest := rl.limits[userID][0]
	return oldest.Add(rl.window)
}

func SanitizeInput(input string) string {
	// Remove potentially dangerous characters
	sanitized := regexp.MustCompile(`[<>{}[\]()&|;]`).ReplaceAllString(input, "")
	
	// Limit length
	if len(sanitized) > 2000 {
		sanitized = sanitized[:2000]
	}
	
	return strings.TrimSpace(sanitized)
}

func HasPermission(userPermissions []Permission, required Permission) bool {
	for _, perm := range userPermissions {
		if perm == required || perm == PermissionAdmin {
			return true
		}
	}
	return false
}

func ValidateURL(inputURL string) bool {
	// Basic URL validation
	urlRegex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(inputURL) {
		return false
	}
	
	// Additional validation to prevent internal IP access
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}
	
	// Check for localhost or IP addresses
	host := parsedURL.Hostname()
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return false
	}
	
	// Check if it's an IP address
	if ip := net.ParseIP(host); ip != nil {
		// Block private IP ranges
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return false
		}
	} else {
		// Check if hostname resolves to internal IP
		ips, err := net.LookupIP(host)
		if err != nil {
			// If we can't resolve the hostname, it might be unsafe
			return false
		}
		
		for _, ip := range ips {
			// Block private IP ranges
			if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				return false
			}
		}
	}
	
	return true
}

func IsValidDiscordID(id string) bool {
	// Discord IDs are numeric strings
	idRegex := regexp.MustCompile(`^\d{17,20}$`)
	return idRegex.MatchString(id)
}

func SanitizeFilename(filename string) string {
	// Remove potentially dangerous characters for filenames
	sanitized := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`).ReplaceAllString(filename, "")
	
	// Limit length
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}
	
	// Remove trailing spaces and dots
	sanitized = strings.TrimRight(sanitized, " .")
	
	// Ensure it's not empty
	if sanitized == "" {
		sanitized = "unnamed"
	}
	
	return sanitized
}

func IsValidChannelName(name string) bool {
	// Discord channel names have specific rules
	if len(name) < 1 || len(name) > 100 {
		return false
	}
	
	// Cannot start or end with dash
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return false
	}
	
	// Only allowed characters
	channelRegex := regexp.MustCompile(`^[a-z0-9-_]+$`)
	return channelRegex.MatchString(name)
}

func ContainsProfanity(text string) bool {
	// Simple profanity filter (in a real implementation, this would be more comprehensive)
	profanityList := []string{"badword1", "badword2", "badword3"}
	
	lowerText := strings.ToLower(text)
	for _, word := range profanityList {
		if strings.Contains(lowerText, word) {
			return true
		}
	}
	
	return false
}