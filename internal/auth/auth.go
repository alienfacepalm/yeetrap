package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const (
	tokenFile       = "token.json"
	credentialsFile = "credentials.json"
)

// Authenticator handles YouTube OAuth2 authentication
type Authenticator struct {
	config    *oauth2.Config
	tokenPath string
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator() (*Authenticator, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	credPath := filepath.Join(configDir, credentialsFile)
	b, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %w\nPlease create %s with your OAuth2 credentials from Google Cloud Console", err, credPath)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	return &Authenticator{
		config:    config,
		tokenPath: filepath.Join(configDir, tokenFile),
	}, nil
}

// Authenticate performs the OAuth2 flow with automatic browser opening
func (a *Authenticator) Authenticate() error {
	// Generate a random state for security
	state := fmt.Sprintf("yeetrap-%d", time.Now().Unix())
	authURL := a.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	
	fmt.Println("🔐 Starting YouTube authentication...")
	fmt.Println("📱 Opening browser for Google OAuth2 login...")
	
	// Try to open browser automatically
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("⚠️  Could not open browser automatically: %v\n", err)
		fmt.Printf("🌐 Please open this URL in your browser:\n%v\n\n", authURL)
	} else {
		fmt.Printf("🌐 Browser opened to: %v\n\n", authURL)
	}
	
	fmt.Println("📋 After logging in and granting permissions:")
	fmt.Println("   1. You'll be redirected to a page that may show an error (this is normal)")
	fmt.Println("   2. Copy the authorization code from the URL or page")
	fmt.Println("   3. Paste it below when prompted")
	fmt.Println()
	fmt.Print("🔑 Enter the authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return fmt.Errorf("unable to read authorization code: %w", err)
	}

	// Clean up the auth code (remove any whitespace)
	authCode = strings.TrimSpace(authCode)
	
	fmt.Println("🔄 Exchanging authorization code for access token...")
	
	tok, err := a.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return fmt.Errorf("unable to retrieve token from web: %w", err)
	}

	fmt.Println("💾 Saving authentication token...")
	return a.saveToken(tok)
}

// GetClient returns an authenticated HTTP client
func (a *Authenticator) GetClient() (*http.Client, error) {
	tok, err := a.loadToken()
	if err != nil {
		return nil, fmt.Errorf("unable to load token: %w\nPlease run 'yeetrap auth' first", err)
	}

	return a.config.Client(context.Background(), tok), nil
}

// saveToken saves a token to a file path
func (a *Authenticator) saveToken(token *oauth2.Token) error {
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(a.tokenPath), 0700); err != nil {
		return err
	}

	f, err := os.OpenFile(a.tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}
	defer f.Close()
	
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return fmt.Errorf("unable to encode token: %w", err)
	}
	
	fmt.Printf("✅ Authentication token saved to: %s\n", a.tokenPath)
	return nil
}

// loadToken retrieves a token from a local file
func (a *Authenticator) loadToken() (*oauth2.Token, error) {
	f, err := os.Open(a.tokenPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// getConfigDir returns the configuration directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".yeetrap")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return configDir, nil
}

// openBrowser opens the specified URL in the default web browser
func openBrowser(url string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	
	return cmd.Start()
}


