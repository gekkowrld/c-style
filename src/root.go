/*
Copyright © 2024 Gekko Wrld
*/
package src

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "betty",
	Short: "A coding style by Holberton School",
	Long: `An opinionated coding style by Holberton School.
Inspired by the Linux Kernel Coding Style with modifications.

When in doubt please refer to the ALX Betty Wiki as the authoriative manual.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}
