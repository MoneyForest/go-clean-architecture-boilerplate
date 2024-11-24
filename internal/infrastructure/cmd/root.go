package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/cmd/http"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/cmd/subscriber"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/cmd/task"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	rootCmd.AddCommand(http.HTTPCmd())
	rootCmd.AddCommand(subscriber.SubscriberCmd())
	rootCmd.AddCommand(task.TaskCmd())

	return rootCmd
}
