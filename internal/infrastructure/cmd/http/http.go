package http

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/controller/http"
)

func HTTPCmd() *cobra.Command {
	httpCmd := &cobra.Command{
		Use:   "http",
		Short: "cli http server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	httpCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "running http server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := http.Run(); err != nil {
				log.Fatal(err)
			}
		},
	})
	return httpCmd
}
