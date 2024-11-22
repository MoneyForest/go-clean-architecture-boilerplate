package service

import (
	"context"
	"testing"
	"time"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/internal/domain/model"
	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
)

func TestValidateMatching(t *testing.T) {
	type args struct {
		me      *model.User
		partner *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK: matching is created successfully",
			args: args{
				me: &model.User{
					ID:        uuid.New(),
					Email:     "me@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				partner: &model.User{
					ID:        uuid.New(),
					Email:     "partner@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "NG: me and partner are the same user",
			args: args{
				me: &model.User{
					ID:        uuid.MustParse("019354c2-47f4-7036-84ff-17ed69ff96e0"),
					Email:     "me@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				partner: &model.User{
					ID:        uuid.MustParse("019354c2-47f4-7036-84ff-17ed69ff96e0"),
					Email:     "partner@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "NG: me is nil",
			args: args{
				me: nil,
				partner: &model.User{
					ID:        uuid.MustParse("019354c2-47f4-7036-84ff-17ed69ff96e0"),
					Email:     "partner@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MatchingDomainService{}
			ctx := context.Background()
			err := s.ValidateMatching(ctx, tt.args.me, tt.args.partner)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatching() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
