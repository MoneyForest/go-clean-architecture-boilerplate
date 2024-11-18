package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewUser(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "success: user is created successfully",
			args: args{
				email: "test@example.com",
			},
			want: &User{
				Email: "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.args.email)

			diff := cmp.Diff(
				got, tt.want,
				cmpopts.IgnoreFields(User{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("NewUser() mismatch (-got +want):\n%s", diff)
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
