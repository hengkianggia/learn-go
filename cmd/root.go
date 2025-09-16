package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "learn",
	Short: "A simple Go application CLI",
	Long:  `A simple Go application CLI with authentication features.`, // Deskripsi panjang
	Run: func(cmd *cobra.Command, args []string) {
		// Default action if no subcommand is given
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

