package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCommand)
}

var rootCmd = &cobra.Command{
	Use:   "resman",
	Short: "Resource Management Server",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
