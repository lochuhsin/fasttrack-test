package cmd

import (
	"cli/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var percentile = &cobra.Command{
	Use:     "percentile",
	Aliases: []string{"p", "-p"},
	Short:   "Get percentile with a given username",
	Long:    "Get percentile with a given username",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a username")
			return
		}
		username := args[0]
		message, err := internal.GetPercentile(username)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(message)
	},
}

func init() {
	rootCmd.AddCommand(percentile)
}
