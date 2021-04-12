package gilfoyle

import (
	"fmt"
	"log"
	"os"

	"github.com/dreamvo/gilfoyle/config"

	"github.com/spf13/cobra"
)

var cfg *config.Config
var cfgFile string
var httpPort int
var addr string

var rootCmd = &cobra.Command{
	Use:     "gilfoyle [OPTIONS] [COMMANDS]",
	Short:   "Cloud-native media streaming server",
	Long:    "Gilfoyle is a web application from the Dreamvo project that runs a self-hosted media streaming server.",
	Example: "gilfoyle serve -p 8080 --config ./gilfoyle.yml",
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file path")
}

func initConfig() {
	if cfgFile != "" {
		c, err := config.NewConfig(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		cfg = c
		return
	}

	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg = c
}

// Execute is a function that executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("root command failed to initialize: %s\n", err)
		os.Exit(1)
	}
}
