package cmd

import (
	"fmt"

	"github.com/AlienFacepalm/YeeTrap/internal/auth"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up YeeTrap OAuth2 app with Google Cloud Console",
	Long: `Set up YeeTrap OAuth2 app with Google Cloud Console.

This command will guide you through creating a Google Cloud project,
enabling the YouTube Data API v3, and setting up OAuth2 credentials.

The setup process includes:
1. Creating a Google Cloud project
2. Enabling YouTube Data API v3
3. Creating OAuth2 desktop app credentials
4. Downloading and placing the credentials file
5. Testing the authentication`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 YeeTrap OAuth2 App Setup")
		fmt.Println("============================")
		fmt.Println()
		
		// Check if already set up
		if err := auth.ValidateCredentials(); err == nil {
			fmt.Println("✅ OAuth2 credentials are already set up!")
			fmt.Println("💡 You can now run: yeetrap auth")
			return nil
		}
		
		// Show setup instructions
		auth.PrintSetupInstructions()
		
		// Check if user wants to test
		fmt.Println()
		fmt.Print("🤔 Would you like to test the setup now? (y/N): ")
		var response string
		fmt.Scanln(&response)
		
		if response == "y" || response == "Y" || response == "yes" {
			fmt.Println()
			fmt.Println("🧪 Testing OAuth2 setup...")
			
			if err := auth.ValidateCredentials(); err != nil {
				fmt.Printf("❌ Setup incomplete: %v\n", err)
				fmt.Println("💡 Please complete the setup steps above and try again")
				return err
			}
			
			fmt.Println("✅ Credentials file found and valid!")
			fmt.Println("🎉 Setup complete! You can now run: yeetrap auth")
		} else {
			fmt.Println()
			fmt.Println("📝 Setup instructions saved. Complete the steps above when ready.")
			fmt.Println("💡 Run 'yeetrap setup' again to test your setup")
		}
		
		return nil
	},
}
