package cmd

import (
	"fmt"

	"github.com/AlienFacepalm/YeeTrap/internal/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with YouTube",
	Long:  `Authenticate with YouTube using OAuth2. This will open a browser window for you to log in.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting YouTube authentication...")
		
		authenticator, err := auth.NewAuthenticator()
		if err != nil {
			return fmt.Errorf("failed to create authenticator: %w", err)
		}

		if err := authenticator.Authenticate(); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		fmt.Println("✓ Authentication successful! Token saved.")
		return nil
	},
}


