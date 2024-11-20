package main

import (
	"os"

	"github.com/MoneyForest/go-clean-boilerplate/internal/infrastructure/cmd"
	_ "github.com/MoneyForest/go-clean-boilerplate/pkg/swag"
)

// @title Go Clean Architecture API
// @version 1.0
// @description This is a sample server using clean architecture.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
