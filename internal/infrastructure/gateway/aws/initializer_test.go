package aws

import (
	"context"
	"testing"
)

func TestLoadAWSConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     AWSConfig
		wantErr bool
	}{
		{
			name: "OK",
			cfg: AWSConfig{
				Environment: "local",
				Region:      "us-west-2",
				Endpoint:    "http://localhost:4566",
			},
			wantErr: false,
		},
		{
			name: "NG_MissingEndpoint",
			cfg: AWSConfig{
				Environment: "local",
				Region:      "us-west-2",
				Endpoint:    "",
			},
			wantErr: true,
		},
		{
			name: "NG_InvalidEnvironment",
			cfg: AWSConfig{
				Environment: "production",
				Region:      "us-west-2",
				Endpoint:    "http://localhost:4566",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := LoadAWSConfig(ctx, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAWSConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
