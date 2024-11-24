package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MoneyForest/go-clean-architecture-boilerplate/pkg/uuid"
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
		wantErr  bool
	}{
		{
			name: "OK: valid matching",
			matching: &Matching{
				ID:        uuid.New(),
				MeID:      meID,
				PartnerID: partnerID,
				Status:    MatchingStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "NG: empty meID",
			matching: &Matching{
				ID:        uuid.New(),
				MeID:      uuid.Nil(),
				PartnerID: partnerID,
				Status:    MatchingStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "NG: empty partnerID",
			matching: &Matching{
				ID:        uuid.New(),
				MeID:      meID,
				PartnerID: uuid.Nil(),
				Status:    MatchingStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "NG: empty status",
			matching: &Matching{
				ID:        uuid.New(),
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "NG: invalid status",
			matching: &Matching{
				ID:        uuid.New(),
				MeID:      meID,
				PartnerID: partnerID,
				Status:    "invalid",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.matching.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatching() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMatching_Accept(t *testing.T) {
	tests := []struct {
		name     string
		matching *Matching
		wantErr  bool
	}{
		{
			name: "OK: matching status is pending",
			matching: &Matching{
				Status: MatchingStatusPending,
			},
			wantErr: false,
		},
		{
			name: "NG: matching status is not pending",
			matching: &Matching{
				Status: MatchingStatusAccepted,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.matching.Accept()
			if (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMatching_Reject(t *testing.T) {
	tests := []struct {
		name     string
		matching *Matching
		wantErr  bool
	}{
		{
			name: "OK: matching status is pending",
			matching: &Matching{
				Status: MatchingStatusPending,
			},
			wantErr: false,
		},
		{
			name: "NG: matching status is not pending",
			matching: &Matching{
				Status: MatchingStatusRejected,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.matching.Reject()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
