package task

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/task"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/task/sample"
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
			if err := task.Run(sample.Run, args); err != nil {
				log.Fatal(err)
			}
		},
	})

	return taskCmd
}
