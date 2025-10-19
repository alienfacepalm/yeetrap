package cmd

import (
	"fmt"
	"os"

	"github.com/AlienFacepalm/YeeTrap/internal/auth"
	"github.com/AlienFacepalm/YeeTrap/internal/downloader"
	"github.com/AlienFacepalm/YeeTrap/internal/youtube"
	"github.com/spf13/cobra"
)

var (
	downloadChannelID string
	downloadMaxVideos int64
	outputDir         string
	quality           string
	concurrent        int
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download videos from a YouTube channel",
	Long:  `Download all videos from your authenticated YouTube channel for backup purposes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

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

		videos, err := ytService.ListChannelVideos(downloadChannelID, downloadMaxVideos)
		if err != nil {
			return fmt.Errorf("failed to list videos: %w", err)
		}

		fmt.Printf("Found %d videos to download\n\n", len(videos))

		dl, err := downloader.NewDownloader(outputDir, quality, concurrent)
		if err != nil {
			return fmt.Errorf("failed to create downloader: %w", err)
		}
		
		if err := dl.DownloadVideos(videos); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		fmt.Println("\n✓ All downloads completed!")
		return nil
	},
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadChannelID, "channel", "c", "", "YouTube channel ID (leave empty to use authenticated user's channel)")
	downloadCmd.Flags().Int64VarP(&downloadMaxVideos, "max", "m", 50, "Maximum number of videos to download")
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", "./downloads", "Output directory for downloaded videos")
	downloadCmd.Flags().StringVarP(&quality, "quality", "q", "best", "Video quality (best, 1080p, 720p, 480p)")
	downloadCmd.Flags().IntVarP(&concurrent, "concurrent", "j", 3, "Number of concurrent downloads")
}


