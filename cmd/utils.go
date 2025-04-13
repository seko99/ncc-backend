package cmd

import (
	"github.com/spf13/cobra"
)

var utilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Utils",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
