package main

import (
	"os"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
