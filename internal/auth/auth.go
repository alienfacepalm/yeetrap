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
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/AlienFacepalm/YeeTrap/internal/constants"
	"github.com/AlienFacepalm/YeeTrap/internal/errors"
	"github.com/AlienFacepalm/YeeTrap/internal/logger"
	"github.com/AlienFacepalm/YeeTrap/internal/retry"
	"github.com/AlienFacepalm/YeeTrap/internal/validation"
)

// Use constants from the constants package

// Authenticator handles YouTube OAuth2 authentication
type Authenticator struct {
	config    *oauth2.Config
	tokenPath string
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator() (*Authenticator, error) {
	logger.Debug("Creating new authenticator")
	
	credPath, err := constants.GetCredentialsPath()
	if err != nil {
		return nil, errors.WrapConfig(err, "failed to get credentials path")
	}

	b, err := os.ReadFile(credPath)
	if err != nil {
		return nil, errors.WrapConfig(err, "unable to read credentials file").
			WithDetails(fmt.Sprintf("Please create %s with your OAuth2 credentials from Google Cloud Console", credPath))
	}

	config, err := google.ConfigFromJSON(b, constants.YouTubeReadonlyScope)
	if err != nil {
		return nil, errors.WrapConfig(err, "unable to parse client secret file to config")
	}

	tokenPath, err := constants.GetTokenPath()
	if err != nil {
		return nil, errors.WrapConfig(err, "failed to get token path")
	}

	logger.Debug("Authenticator created successfully")
	return &Authenticator{
		config:    config,
		tokenPath: tokenPath,
	}, nil
}

// Authenticate performs the OAuth2 flow with automatic browser opening
func (a *Authenticator) Authenticate() error {
	logger.Info("Starting YouTube authentication")
	
	// Generate a random state for security
	state := fmt.Sprintf("yeetrap-%d", time.Now().Unix())
	authURL := a.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	
	fmt.Println("🔐 Starting YouTube authentication...")
	fmt.Println("📱 Opening browser for Google OAuth2 login...")
	
	// Try to open browser automatically
	if err := openBrowser(authURL); err != nil {
		logger.Warn("Could not open browser automatically: %v", err)
		fmt.Printf("⚠️  Could not open browser automatically: %v\n", err)
		fmt.Printf("🌐 Please open this URL in your browser:\n%v\n\n", authURL)
	} else {
		logger.Info("Browser opened successfully")
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
		return errors.WrapAuth(err, "unable to read authorization code")
	}

	// Validate and clean up the auth code
	authCode, err := validation.ValidateAndSanitizeInput(authCode, "auth_code")
	if err != nil {
		return err
	}
	
	fmt.Println("🔄 Exchanging authorization code for access token...")
	
	// Use retry logic for token exchange
	var tok *oauth2.Token
	err = retry.RetryAPIOperation(func() error {
		var err error
		tok, err = a.config.Exchange(context.TODO(), authCode)
		return err
	})
	
	if err != nil {
		return errors.WrapAuth(err, "unable to retrieve token from web")
	}

	fmt.Println("💾 Saving authentication token...")
	return a.saveToken(tok)
}

// GetClient returns an authenticated HTTP client
func (a *Authenticator) GetClient() (*http.Client, error) {
	logger.Debug("Getting authenticated HTTP client")
	
	tok, err := a.loadToken()
	if err != nil {
		return nil, errors.WrapAuth(err, "unable to load token").
			WithDetails("Please run 'yeetrap auth' first")
	}

	logger.Debug("HTTP client created successfully")
	return a.config.Client(context.Background(), tok), nil
}

// saveToken saves a token to a file path
func (a *Authenticator) saveToken(token *oauth2.Token) error {
	logger.Debug("Saving authentication token")
	
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(a.tokenPath), 0700); err != nil {
		return errors.WrapFile(err, "unable to create token directory")
	}

	f, err := os.OpenFile(a.tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.WrapFile(err, "unable to create token file")
	}
	defer f.Close()
	
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return errors.WrapFile(err, "unable to encode token")
	}
	
	logger.Info("Authentication token saved successfully")
	fmt.Printf("✅ Authentication token saved to: %s\n", a.tokenPath)
	return nil
}

// loadToken retrieves a token from a local file
func (a *Authenticator) loadToken() (*oauth2.Token, error) {
	logger.Debug("Loading authentication token")
	
	f, err := os.Open(a.tokenPath)
	if err != nil {
		return nil, errors.WrapFile(err, "unable to open token file")
	}
	defer f.Close()
	
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, errors.WrapFile(err, "unable to decode token")
	}
	
	logger.Debug("Authentication token loaded successfully")
	return tok, nil
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


