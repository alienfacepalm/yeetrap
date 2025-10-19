package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yeetrap",
	Short: "YeeTrap - YouTube Video Downloader for Content Creators",
	Long: `YeeTrap is a tool for YouTube content creators to download their videos in bulk for backup purposes.
It authenticates with YouTube and allows you to download all videos from your channel.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(versionCmd)
}

