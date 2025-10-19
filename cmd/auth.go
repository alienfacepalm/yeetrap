package cmd

import (
	"fmt"

	"github.com/AlienFacepalm/YeeTrap/internal/auth"
	"github.com/spf13/cobra"
)

var (
	showSetup bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with YouTube using OAuth2",
	Long: `Authenticate with YouTube using OAuth2. This will open a browser window for you to log in.

YeeTrap uses Google OAuth2 to securely access your YouTube channel data.
The authentication token is saved locally and reused for future sessions.

If you haven't set up OAuth2 credentials yet, use: yeetrap auth --setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showSetup {
			auth.PrintSetupInstructions()
			return nil
		}

		fmt.Println("🔐 YeeTrap YouTube Authentication")
		fmt.Println("=================================")
		fmt.Println()
		
		// Check if credentials exist
		if err := auth.ValidateCredentials(); err != nil {
			fmt.Printf("❌ %v\n\n", err)
			fmt.Println("💡 Run 'yeetrap auth --setup' for detailed setup instructions")
			return err
		}
		
		authenticator, err := auth.NewAuthenticator()
		if err != nil {
			return fmt.Errorf("failed to create authenticator: %w", err)
		}

		if err := authenticator.Authenticate(); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		fmt.Println()
		fmt.Println("🎉 Authentication successful!")
		fmt.Println("✅ You can now use 'yeetrap list' and 'yeetrap download' commands")
		fmt.Println("💡 Your authentication token is saved and will be reused automatically")
		return nil
	},
}

func init() {
	authCmd.Flags().BoolVar(&showSetup, "setup", false, "Show detailed OAuth2 setup instructions")
}


