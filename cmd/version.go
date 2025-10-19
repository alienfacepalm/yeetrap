package cmd

import (
	"fmt"

	"github.com/AlienFacepalm/YeeTrap/internal/constants"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of YeeTrap",
	Long:  `Print the version number of YeeTrap`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", constants.AppName, constants.AppVersion)
		fmt.Println(constants.AppDescription)
		fmt.Println(constants.AppURL)
	},
}

