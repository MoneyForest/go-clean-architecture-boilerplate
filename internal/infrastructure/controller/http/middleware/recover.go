package middleware

import (
	"net/http"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/error"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure/controller/http/response"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				appErr := error.NewDomainError(
					error.Critical,
					"Internal server error",
					nil,
					map[string]interface{}{"panic": err},
				)
				response.WriteError(w, appErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
