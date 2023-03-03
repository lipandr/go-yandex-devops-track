package config

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestNewServer tests os env vars are parsed correctly.
func TestNewServer(t *testing.T) {
	tests := []struct {
		name  string
		input Config
		want  Config
	}{
		{
			name: "PollInterval is 1 second, Report interval is 2 seconds",
			input: Config{
				Address:       "localhost:8888",
				StoreInterval: 120 * time.Second,
				StoreFile:     "/test.tmp",
				Restore:       false,
			},
			want: Config{
				Address:       "localhost:8888",
				StoreInterval: 120 * time.Second,
				StoreFile:     "/test.tmp",
				Restore:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.Setenv("ADDRESS", tt.input.Address))
			require.NoError(t, os.Setenv("STORE_INTERVAL", tt.input.StoreInterval.String()))
			require.NoError(t, os.Setenv("STORE_FILE", tt.input.StoreFile))
			require.NoError(t, os.Setenv("RESTORE", strconv.FormatBool(tt.input.Restore)))
			got := NewServer()

			require.Equalf(t, tt.want.Address, got.Address,
				"NewServer() want: %v, got: %v", tt.want.Address, got.Address)
			require.Equalf(t, tt.want.StoreInterval, got.StoreInterval,
				"NewServer() want: %v, got: %v", tt.want.StoreInterval, got.StoreInterval)
			require.Equalf(t, tt.want.StoreFile, got.StoreFile,
				"NewServer() want: %v, got: %v", tt.want.StoreFile, got.StoreFile)
			require.Equalf(t, tt.want.Restore, got.Restore,
				"NewServer() want: %v, got: %v", tt.want.Restore, got.Restore)
		})
	}
}
