package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlienFacepalm/YeeTrap/internal/constants"
	"github.com/AlienFacepalm/YeeTrap/internal/errors"
)

// YouTube channel ID pattern
var channelIDPattern = regexp.MustCompile(`^UC[a-zA-Z0-9_-]{22}$`)

// YouTube video ID pattern
var videoIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

// ValidateChannelID validates a YouTube channel ID
func ValidateChannelID(channelID string) error {
	if channelID == "" {
		return nil // Empty is valid (means use authenticated user's channel)
	}
	
	if !channelIDPattern.MatchString(channelID) {
		return errors.NewValidationError(fmt.Sprintf("invalid channel ID format: %s", channelID)).
			WithDetails("Channel ID should start with 'UC' followed by 22 characters")
	}
	
	return nil
}

// ValidateVideoID validates a YouTube video ID
func ValidateVideoID(videoID string) error {
	if videoID == "" {
		return errors.NewValidationError("video ID cannot be empty")
	}
	
	if !videoIDPattern.MatchString(videoID) {
		return errors.NewValidationError(fmt.Sprintf("invalid video ID format: %s", videoID)).
			WithDetails("Video ID should be 11 characters long")
	}
	
	return nil
}

// ValidateQuality validates video quality parameter
func ValidateQuality(quality string) error {
	if quality == "" {
		return errors.NewValidationError("quality cannot be empty")
	}
	
	for _, validQuality := range constants.SupportedQualities {
		if quality == validQuality {
			return nil
		}
	}
	
	return errors.NewValidationError(fmt.Sprintf("invalid quality: %s", quality)).
		WithDetails(fmt.Sprintf("Supported qualities: %s", strings.Join(constants.SupportedQualities, ", ")))
}

// ValidateConcurrency validates concurrent download count
func ValidateConcurrency(concurrent int) error {
	if concurrent < 1 {
		return errors.NewValidationError("concurrency must be at least 1")
	}
	
	if concurrent > 10 {
		return errors.NewValidationError("concurrency should not exceed 10").
			WithDetails("High concurrency may cause rate limiting or system issues")
	}
	
	return nil
}

// ValidateMaxVideos validates maximum video count
func ValidateMaxVideos(maxVideos int64) error {
	if maxVideos < 0 {
		return errors.NewValidationError("max videos cannot be negative")
	}
	
	if maxVideos > 10000 {
		return errors.NewValidationError("max videos should not exceed 10000").
			WithDetails("Large numbers may cause API quota issues")
	}
	
	return nil
}

// ValidateOutputDir validates output directory path
func ValidateOutputDir(outputDir string) error {
	if outputDir == "" {
		return errors.NewValidationError("output directory cannot be empty")
	}
	
	// Check for invalid characters in path
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(outputDir, char) {
			return errors.NewValidationError(fmt.Sprintf("output directory contains invalid character: %s", char))
		}
	}
	
	return nil
}

// ValidateAuthCode validates OAuth2 authorization code
func ValidateAuthCode(authCode string) error {
	if authCode == "" {
		return errors.NewValidationError("authorization code cannot be empty")
	}
	
	// Basic validation - auth codes are typically longer and contain various characters
	if len(authCode) < 10 {
		return errors.NewValidationError("authorization code appears to be too short")
	}
	
	return nil
}

// SanitizeFilename sanitizes a filename for safe filesystem usage
func SanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename
	
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	// Remove leading/trailing whitespace and dots
	result = strings.Trim(result, " .")
	
	// Limit length to prevent filesystem issues
	if len(result) > 200 {
		result = result[:200]
	}
	
	// Ensure it's not empty
	if result == "" {
		result = "untitled"
	}
	
	return result
}

// ValidateAndSanitizeInput validates and sanitizes user input
func ValidateAndSanitizeInput(input string, inputType string) (string, error) {
	if input == "" {
		return "", errors.NewValidationError(fmt.Sprintf("%s cannot be empty", inputType))
	}
	
	// Basic sanitization
	sanitized := strings.TrimSpace(input)
	
	// Additional validation based on input type
	switch inputType {
	case "channel_id":
		if err := ValidateChannelID(sanitized); err != nil {
			return "", err
		}
	case "video_id":
		if err := ValidateVideoID(sanitized); err != nil {
			return "", err
		}
	case "quality":
		if err := ValidateQuality(sanitized); err != nil {
			return "", err
		}
	case "auth_code":
		if err := ValidateAuthCode(sanitized); err != nil {
			return "", err
		}
	}
	
	return sanitized, nil
}

// ValidateConfig validates application configuration
func ValidateConfig(config interface{}) error {
	// This would be expanded based on the actual config structure
	// For now, it's a placeholder for future configuration validation
	return nil
}
