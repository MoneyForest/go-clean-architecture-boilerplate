package validate

import (
	"testing"

	e "github.com/MoneyForest/go-clean-boilerplate/internal/domain/error"
	"github.com/MoneyForest/go-clean-boilerplate/internal/domain/model"
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		wantErr error
	}{
		{
			name: "success: valid user with email",
			user: &model.User{
				Email: "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "error: empty email",
			user: &model.User{
				Email: "",
			},
			wantErr: e.ErrUserEmailIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUser(tt.user)
			if err != tt.wantErr {
				t.Errorf("ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
