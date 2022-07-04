package config

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TestGetConfig load default config from .env.example
func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{
			name: "success",
			want: Config{
				BaseURL:        "http://localhost:8100",
				ServiceName:    "api",
				Port:           "8100",
				AllowedOrigins: "https://*.dwarves.foundation/;https://*.vercel.app",
				AccessTokenTTL: time.Hour * 24 * time.Duration(30),
				Env:            "dev",
				JWTSecret:      []byte("sample_secret"),
				Version:        "1.0.0",
			},
		},
	}

	ignoreFields := cmpopts.IgnoreFields(Config{},
		"DBHost", "DBPort", "DBUser", "DBName", "DBPass", "DBSSLMode",
	)

	// override env file
	envDir = "../.."
	envFileName = ".env.sample"
	defer func() {
		envDir = "."
		envFileName = ".env"
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfig()
			if !cmp.Equal(got, tt.want, ignoreFields) {
				t.Errorf("GetConfig() = %v, want %v \n diff: %v", got, tt.want, cmp.Diff(got, tt.want, ignoreFields))
				t.FailNow()
			}
		})
	}
}
