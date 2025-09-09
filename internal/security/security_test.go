package security

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(5, time.Minute)

	if rl.maxReq != 5 {
		t.Errorf("Expected maxReq to be 5, got %d", rl.maxReq)
	}

	if rl.window != time.Minute {
		t.Errorf("Expected window to be 1 minute, got %v", rl.window)
	}

	if rl.limits == nil {
		t.Error("Expected limits map to be initialized, got nil")
	}
}

func TestRateLimiterIsAllowed(t *testing.T) {
	rl := NewRateLimiter(2, time.Second)

	userID := "test_user"

	// First request should be allowed
	if !rl.IsAllowed(userID) {
		t.Error("First request should be allowed")
	}

	// Second request should be allowed
	if !rl.IsAllowed(userID) {
		t.Error("Second request should be allowed")
	}

	// Third request should be denied
	if rl.IsAllowed(userID) {
		t.Error("Third request should be denied")
	}
}

func TestSanitizeInput(t *testing.T) {
	// Test removing dangerous characters
	input := "<script>alert('xss')</script>Hello & welcome;"
	expected := "Hello  welcome"
	result := SanitizeInput(input)

	if result != expected {
		t.Errorf("Expected sanitized input to be '%s', got '%s'", expected, result)
	}

	// Test length limiting
	longInput := string(make([]byte, 2500))
	result = SanitizeInput(longInput)

	if len(result) > 2000 {
		t.Errorf("Expected result to be limited to 2000 characters, got %d", len(result))
	}
}

func TestHasPermission(t *testing.T) {
	userPerms := []Permission{PermissionMusic, PermissionDownload}

	// Test having specific permission
	if !HasPermission(userPerms, PermissionMusic) {
		t.Error("User should have music permission")
	}

	// Test not having specific permission
	if HasPermission(userPerms, PermissionAI) {
		t.Error("User should not have AI permission")
	}

	// Test admin permission granting access to all
	adminPerms := []Permission{PermissionAdmin}
	if !HasPermission(adminPerms, PermissionAI) {
		t.Error("Admin should have access to all permissions")
	}
}

func TestValidateURL(t *testing.T) {
	// Test valid URL
	validURL := "https://example.com"
	if !ValidateURL(validURL) {
		t.Errorf("Expected URL '%s' to be valid", validURL)
	}

	// Test invalid URL
	invalidURL := "not a url"
	if ValidateURL(invalidURL) {
		t.Errorf("Expected URL '%s' to be invalid", invalidURL)
	}

	// Test localhost URL (should be invalid)
	localhostURL := "http://localhost:8080"
	if ValidateURL(localhostURL) {
		t.Errorf("Expected localhost URL '%s' to be invalid", localhostURL)
	}
}

func TestIsValidDiscordID(t *testing.T) {
	// Test valid Discord ID
	validID := "123456789012345678"
	if !IsValidDiscordID(validID) {
		t.Errorf("Expected ID '%s' to be valid", validID)
	}

	// Test invalid Discord ID (too short)
	invalidID := "12345"
	if IsValidDiscordID(invalidID) {
		t.Errorf("Expected ID '%s' to be invalid", invalidID)
	}

	// Test invalid Discord ID (contains letters)
	invalidID2 := "123456789012345abc"
	if IsValidDiscordID(invalidID2) {
		t.Errorf("Expected ID '%s' to be invalid", invalidID2)
	}
}