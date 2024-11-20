package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MoneyForest/go-clean-boilerplate/pkg/uuid"
)

func TestNewMatching(t *testing.T) {
	type args struct {
		params InputMatchingParams
	}
	meID := uuid.New()
	partnerID := uuid.New()

	tests := []struct {
		name string
		args args
		want *Matching
	}{
		{
			name: "success: matching is created successfully",
			args: args{
				params: InputMatchingParams{
					MeID:      meID,
					PartnerID: partnerID,
					Status:    "test",
				},
			},
			want: &Matching{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMatching(tt.args.params)

			diff := cmp.Diff(
				got, tt.want,
				cmpopts.IgnoreFields(Matching{}, "ID", "CreatedAt", "UpdatedAt"),
			)
			if diff != "" {
				t.Errorf("NewMatching() mismatching (-got +want):\n%s", diff)
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

func TestValidateMatching(t *testing.T) {
	meID := uuid.New()
	partnerID := uuid.New()

	tests := []struct {
		name     string
		matching *Matching
		wantErr  error
	}{
		{
			name: "success: valid matching",
			matching: &Matching{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    MatchingStatusPending,
			},
			wantErr: nil,
		},
		{
			name: "error: empty meID",
			matching: &Matching{
				MeID:      uuid.Nil(),
				PartnerID: partnerID,
				Status:    MatchingStatusPending,
			},
			wantErr: ErrMatchingMeOrPartnerIDIsRequired,
		},
		{
			name: "error: empty partnerID",
			matching: &Matching{
				MeID:      meID,
				PartnerID: uuid.Nil(),
				Status:    MatchingStatusPending,
			},
			wantErr: ErrMatchingMeOrPartnerIDIsRequired,
		},
		{
			name: "error: empty status",
			matching: &Matching{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "",
			},
			wantErr: ErrMatchingStatusIsRequired,
		},
		{
			name: "error: invalid status",
			matching: &Matching{
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "invalid",
			},
			wantErr: ErrMatchingStatusIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.matching.Validate()
			if err != tt.wantErr {
				t.Errorf("ValidateMatching() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
