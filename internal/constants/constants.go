package constants

import (
	"os"
	"path/filepath"
)

// Application constants
const (
	AppName        = "YeeTrap"
	AppVersion     = "1.0.0"
	AppDescription = "YouTube Video Downloader for Content Creators"
	AppURL         = "https://github.com/AlienFacepalm/YeeTrap"
)

// File and directory constants
const (
	ConfigDirName     = ".yeetrap"
	CredentialsFile   = "credentials.json"
	TokenFile         = "token.json"
	ConfigFile        = "config.json"
	DefaultOutputDir  = "./downloads"
)

// YouTube API constants
const (
	YouTubeReadonlyScope = "https://www.googleapis.com/auth/youtube.readonly"
	MaxVideosPerPage     = 50
	DefaultMaxVideos     = 50
	DefaultConcurrency   = 3
)

// Video quality options
const (
	QualityBest  = "best"
	Quality1080p = "1080p"
	Quality720p  = "720p"
	Quality480p  = "480p"
)

// Supported video qualities
var SupportedQualities = []string{QualityBest, Quality1080p, Quality720p, Quality480p}

// Default configuration values
const (
	DefaultQuality      = QualityBest
	DefaultChannelID    = ""
	DefaultMaxConcurrent = 3
)

// Error messages
const (
	ErrCredentialsNotFound = "credentials file not found"
	ErrTokenNotFound      = "authentication token not found"
	ErrYtDlpNotFound      = "yt-dlp is not installed or not in PATH"
	ErrInvalidQuality     = "invalid video quality specified"
	ErrInvalidChannelID   = "invalid channel ID format"
	ErrAPIQuotaExceeded   = "YouTube API quota exceeded"
	ErrNetworkError       = "network error occurred"
	ErrFileSystemError    = "file system error occurred"
)

// Success messages
const (
	MsgAuthSuccess      = "Authentication successful!"
	MsgDownloadComplete = "All downloads completed!"
	MsgSetupComplete    = "Setup complete!"
	MsgTokenSaved       = "Authentication token saved"
)

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetCredentialsPath returns the full path to the credentials file
func GetCredentialsPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, CredentialsFile), nil
}

// GetTokenPath returns the full path to the token file
func GetTokenPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, TokenFile), nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, ConfigFile), nil
}
