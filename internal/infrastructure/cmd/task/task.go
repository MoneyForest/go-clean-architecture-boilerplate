package task

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/task"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/task/enqueue_user_deletion"
)

func TaskCmd() *cobra.Command {
	taskCmd := &cobra.Command{
		Use: "task",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	taskCmd.AddCommand(&cobra.Command{
		Use:   "sample",
		Short: "sample",
		Run: func(cmd *cobra.Command, args []string) {
			if err := task.Run(enqueue_user_deletion.Run, args); err != nil {
				log.Fatal(err)
			}
		},
	})

	return taskCmd
}
