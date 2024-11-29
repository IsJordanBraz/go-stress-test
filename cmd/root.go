package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

var rootCmd = &cobra.Command{
	Use:   "go-stress-test",
	Short: "A brief description of your application",
	Long:  "A brief description of your application",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
