package subscriber

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/subscriber"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/subscriber/dequeue_and_delete_user"
)

func SubscriberCmd() *cobra.Command {
	subscriberCmd := &cobra.Command{
		Use: "subscriber",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	subscriberCmd.AddCommand(&cobra.Command{
		Use:   "sample",
		Short: "sample",
		Run: func(cmd *cobra.Command, args []string) {
			if err := subscriber.Run(dequeue_and_delete_user.Run, args); err != nil {
				log.Fatal(err)
			}
		},
	})

	return subscriberCmd
}
