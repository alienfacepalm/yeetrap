package cmd

import (
	"fmt"

	"github.com/AlienFacepalm/YeeTrap/internal/auth"
	"github.com/AlienFacepalm/YeeTrap/internal/youtube"
	"github.com/spf13/cobra"
)

var (
	channelID string
	maxVideos int64
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List videos from a YouTube channel",
	Long:  `List all videos from your authenticated YouTube channel.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		authenticator, err := auth.NewAuthenticator()
		if err != nil {
			return fmt.Errorf("failed to create authenticator: %w", err)
		}

		client, err := authenticator.GetClient()
		if err != nil {
			return fmt.Errorf("failed to get authenticated client: %w", err)
		}

		ytService, err := youtube.NewService(client)
		if err != nil {
			return fmt.Errorf("failed to create YouTube service: %w", err)
		}

		videos, err := ytService.ListChannelVideos(channelID, maxVideos)
		if err != nil {
			return fmt.Errorf("failed to list videos: %w", err)
		}

		fmt.Printf("Found %d videos:\n\n", len(videos))
		for i, video := range videos {
			fmt.Printf("%d. %s\n", i+1, video.Title)
			fmt.Printf("   ID: %s\n", video.ID)
			fmt.Printf("   URL: https://www.youtube.com/watch?v=%s\n", video.ID)
			fmt.Printf("   Published: %s\n\n", video.PublishedAt)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&channelID, "channel", "c", "", "YouTube channel ID (leave empty to use authenticated user's channel)")
	listCmd.Flags().Int64VarP(&maxVideos, "max", "m", 50, "Maximum number of videos to list")
}


