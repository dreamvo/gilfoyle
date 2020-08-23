package cmd

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

func init() {
	// Register command
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print configuration",
	Run: func(cmd *cobra.Command, args []string) {
		d, err := yaml.Marshal(config.GetConfig())
		if err != nil {
			panic(err)
		}
		fmt.Println(string(d))
	},
}
