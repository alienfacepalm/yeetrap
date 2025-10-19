package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of YeeTrap",
	Long:  `Print the version number of YeeTrap`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("YeeTrap v%s\n", version)
		fmt.Println("YouTube Video Downloader for Content Creators")
		fmt.Println("https://github.com/AlienFacepalm/YeeTrap")
	},
}

