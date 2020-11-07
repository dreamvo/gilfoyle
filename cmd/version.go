package cmd

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/spf13/cobra"
)

func init() {
	// Register command
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print program version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Gilfoyle version %s commit %s\n", config.Version, config.Commit)
	},
}
