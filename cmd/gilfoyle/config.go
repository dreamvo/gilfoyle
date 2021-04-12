package gilfoyle

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	// Register command
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print configuration",
	Run: func(cmd *cobra.Command, args []string) {
		d, err := yaml.Marshal(cfg)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(d))
	},
}
