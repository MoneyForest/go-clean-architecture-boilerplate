package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

func TestNewMatch(t *testing.T) {
	type args struct {
		params InputMatchParams
	}
	meID := uuid.New()
	partnerID := uuid.New()

	tests := []struct {
		name string
		args args
		want *Match
	}{
		{
			name: "success: match is created successfully",
			args: args{
				params: InputMatchParams{
					MeID:      meID,
					PartnerID: partnerID,
					Status:    "test",
				},
			},
			want: &Match{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMatch(tt.args.params)

			diff := cmp.Diff(
				got, tt.want,
				cmpopts.IgnoreFields(Match{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("NewMatch() mismatch (-got +want):\n%s", diff)
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

func TestValidateMatch(t *testing.T) {
	meID := uuid.New()
	partnerID := uuid.New()

	tests := []struct {
		name    string
		match   *Match
		wantErr error
	}{
		{
			name: "success: valid match",
			match: &Match{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    MatchStatusPending,
			},
			wantErr: nil,
		},
		{
			name: "error: empty meID",
			match: &Match{
				MeID:      uuid.Nil(),
				PartnerID: partnerID,
				Status:    MatchStatusPending,
			},
			wantErr: ErrMatchMeOrPartnerIDIsRequired,
		},
		{
			name: "error: empty partnerID",
			match: &Match{
				MeID:      meID,
				PartnerID: uuid.Nil(),
				Status:    MatchStatusPending,
			},
			wantErr: ErrMatchMeOrPartnerIDIsRequired,
		},
		{
			name: "error: empty status",
			match: &Match{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "",
			},
			wantErr: ErrMatchStatusIsRequired,
		},
		{
			name: "error: invalid status",
			match: &Match{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "invalid",
			},
			wantErr: ErrMatchStatusIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.match.Validate()
			if err != tt.wantErr {
				t.Errorf("ValidateMatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
