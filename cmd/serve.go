package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var httpPort int

func init() {
	// Register command
	rootCmd.AddCommand(serveCmd)

	// Register flags
	serveCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 3000, "HTTP port")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve REST API",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve")
	},
}
