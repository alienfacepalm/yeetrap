package auth

import (
	"fmt"
	"os"

	"github.com/AlienFacepalm/YeeTrap/internal/constants"
	"github.com/AlienFacepalm/YeeTrap/internal/errors"
	"github.com/AlienFacepalm/YeeTrap/internal/logger"
)

// AppConfig represents the OAuth2 app configuration
type AppConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ProjectID    string `json:"project_id"`
	AppName      string `json:"app_name"`
	RedirectURI  string `json:"redirect_uri"`
}

// GetAppInfo returns information about the OAuth2 app setup
func GetAppInfo() *AppConfig {
	return &AppConfig{
		AppName:     constants.AppName,
		RedirectURI: "http://localhost:8080/callback",
	}
}

// ValidateCredentials checks if the credentials file exists and is valid
func ValidateCredentials() error {
	logger.Debug("Validating credentials file")
	
	credPath, err := constants.GetCredentialsPath()
	if err != nil {
		return errors.WrapConfig(err, "failed to get credentials path")
	}

	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		return errors.NewConfigError("credentials file not found").
			WithDetails(fmt.Sprintf("Expected location: %s", credPath))
	}

	// Try to read and parse the credentials
	_, err = os.ReadFile(credPath)
	if err != nil {
		return errors.WrapFile(err, "unable to read credentials file")
	}

	logger.Debug("Credentials file validation successful")
	return nil
}

// PrintSetupInstructions prints detailed setup instructions
func PrintSetupInstructions() {
	logger.Info("Printing setup instructions")
	
	credPath, _ := constants.GetCredentialsPath()
	
	fmt.Println("🔧 YeeTrap OAuth2 App Setup Instructions")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("1. 🌐 Go to Google Cloud Console:")
	fmt.Println("   https://console.cloud.google.com/")
	fmt.Println()
	fmt.Println("2. 📁 Create a new project:")
	fmt.Println("   - Click 'Select a project' → 'New Project'")
	fmt.Println("   - Name it 'YeeTrap' or similar")
	fmt.Println("   - Click 'Create'")
	fmt.Println()
	fmt.Println("3. 🔌 Enable YouTube Data API v3:")
	fmt.Println("   - Go to 'APIs & Services' → 'Library'")
	fmt.Println("   - Search for 'YouTube Data API v3'")
	fmt.Println("   - Click on it and press 'Enable'")
	fmt.Println()
	fmt.Println("4. 🔐 Create OAuth2 credentials:")
	fmt.Println("   - Go to 'APIs & Services' → 'Credentials'")
	fmt.Println("   - Click 'Create Credentials' → 'OAuth client ID'")
	fmt.Println("   - If prompted, configure OAuth consent screen:")
	fmt.Println("     * Choose 'External' user type")
	fmt.Println("     * App name: 'YeeTrap'")
	fmt.Println("     * Support email: your email")
	fmt.Println("     * Add your email as a test user")
	fmt.Println("     * Save and continue")
	fmt.Println("   - Back to 'Create OAuth client ID':")
	fmt.Println("     * Application type: 'Desktop app'")
	fmt.Println("     * Name: 'YeeTrap Desktop'")
	fmt.Println("     * Click 'Create'")
	fmt.Println()
	fmt.Println("5. 💾 Download and place credentials:")
	fmt.Println("   - Download the JSON file")
	fmt.Println("   - Rename it to 'credentials.json'")
	fmt.Println("   - Place it at:")
	fmt.Printf("     %s\n", credPath)
	fmt.Println()
	fmt.Println("6. ✅ Test the setup:")
	fmt.Println("   - Run: yeetrap auth")
	fmt.Println("   - Your browser should open automatically")
	fmt.Println("   - Login with your Google/YouTube account")
	fmt.Println("   - Grant permissions to YeeTrap")
	fmt.Println("   - Copy the authorization code back to terminal")
	fmt.Println()
	fmt.Println("🎉 That's it! You're ready to use YeeTrap!")
}
