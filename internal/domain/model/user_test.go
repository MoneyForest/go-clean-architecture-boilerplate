package model

import (
	"testing"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewUser(t *testing.T) {
	type args struct {
		params InputUserParams
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "success: user is created successfully",
			args: args{
				params: InputUserParams{
					Email: "test@example.com",
				},
			},
			want: &User{
				Email: "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.args.params)

			diff := cmp.Diff(
				got, tt.want,
				cmpopts.IgnoreFields(User{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("NewUser() mismatching (-got +want):\n%s", diff)
			}

			now := time.Now()
			if got.CreatedAt.Sub(now) > time.Second {
				t.Error("CreatedAt should be close to current time")
			}
			if got.UpdatedAt.Sub(now) > time.Second {
				t.Error("UpdatedAt should be close to current time")
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "success: valid user with email",
			user: &User{
				ID:        uuid.New(),
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "error: invalid email",
			user: &User{
				ID:        uuid.New(),
				Email:     "invalid-email",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
