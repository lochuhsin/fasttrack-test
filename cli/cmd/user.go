package cmd

import (
	"cli/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var start = &cobra.Command{
	Use:     "start",
	Aliases: []string{"c", "-c"},
	Short:   "start answering questions",
	Long:    "with one external argument which defines the user name",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username, err := internal.CreateNamePrompt()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = internal.CreateUser(username)
		if err != nil {
			fmt.Println(err)
			return
		}
		qResp, err := internal.GetQuestion()
		if err != nil {
			fmt.Println(err)
			return
		}

		ans, err := internal.RunProblemPrompt(qResp.Questions)
		if err != nil {
			fmt.Println(err)
			return
		}

		str, err := internal.Submit(username, ans)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(str)
	},
}

func init() {
	rootCmd.AddCommand(start)
}
