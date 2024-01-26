/*
Copyright © 2024 Gekko Wrld

*/

package src

import (
	"github.com/spf13/cobra"
)

var FlagsPassed struct {
  Verbose bool
  Quiet bool
}

var styleCmd = &cobra.Command{
	Use:   "style",
	Short: "Check if the code complies with the coding style",

	Run: func(cmd *cobra.Command, args []string) {
    indentation("/home/gekkowrld/code/betty/ninde_test.c")
	},
}

func init() {
	rootCmd.AddCommand(styleCmd)
}
