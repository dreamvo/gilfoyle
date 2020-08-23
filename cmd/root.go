package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "gilfoyle [COMMANDS] [OPTIONS]",
	Short:   "Video streaming API server",
	Long:    "Gilfoyle is a web application from the Dreamvo project that runs a self-hosted video streaming server.",
	Example: "gilfoyle serve -p 8080",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default ./gilfoyle.yaml")
}

// Execute is a function that executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
