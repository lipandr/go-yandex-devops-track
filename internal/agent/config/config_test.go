package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestNewAgent tests os env vars are parsed correctly.
func TestNewAgent(t *testing.T) {
	tests := []struct {
		name  string
		input Config
		want  Config
	}{
		{
			name: "PollInterval is 1 second, Report interval is 2 seconds",
			input: Config{
				Address:        "localhost:8080",
				PollInterval:   1 * time.Second,
				ReportInterval: 2 * time.Second,
			},
			want: Config{
				Address:        "localhost:8080",
				PollInterval:   1 * time.Second,
				ReportInterval: 2 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.Setenv("ADDRESS", tt.input.Address))
			require.NoError(t, os.Setenv("POLL_INTERVAL", tt.input.PollInterval.String()))
			require.NoError(t, os.Setenv("REPORT_INTERVAL", tt.input.ReportInterval.String()))
			got := NewAgent()

			require.Equalf(t, tt.want.Address, got.Address,
				"NewAgent() want: %v, got: %v", tt.want.Address, got.Address)
			require.Equalf(t, tt.want.PollInterval, got.PollInterval,
				"NewAgent() want: %v, got: %v", tt.want.PollInterval, got.PollInterval)
			require.Equalf(t, tt.want.ReportInterval, got.ReportInterval,
				"NewAgent() want: %v, got: %v", tt.want.ReportInterval, got.ReportInterval)

		})
	}
}
