package cmd

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gilfoyle [COMMANDS] [OPTIONS]",
	Short:   "Video streaming API server",
	Long:    "Gilfoyle is a web application from the Dreamvo project that runs a self-hosted video streaming server.",
	Example: "gilfoyle serve -p 8080",
}

func init() {
	cfgFile := rootCmd.PersistentFlags().String("config", "", "config file path (default ./gilfoyle.yaml")

	err := config.NewConfig(cfgFile)
	if err != nil {
		panic(err)
	}
}

// Execute is a function that executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
