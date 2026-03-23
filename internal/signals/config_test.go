package signals

import (
	"strings"
	"testing"
	"time"
)

func TestIdentitySignalsConfig_Validate(t *testing.T) {
	validDebounce := DebouncerConfig{
		MinFrequency: time.Second,
		MaxBulkSize:  100,
	}
	tests := []struct {
		name    string
		cfg     IdentitySignalsConfig
		wantErr string
	}{
		{
			name:    "disabled config is always valid",
			cfg:     IdentitySignalsConfig{Enabled: false},
			wantErr: "",
		},
		{
			name: "enabled without ducklake fails",
			cfg: IdentitySignalsConfig{
				Enabled: true,
				Store: StoreConfig{
					Debounce: validDebounce,
					DuckLake: DuckLakeConfig{Enabled: false},
				},
			},
			wantErr: "Store.DuckLake.Enabled=true",
		},
		{
			name: "enabled with ducklake but no data path fails",
			cfg: IdentitySignalsConfig{
				Enabled: true,
				Store: StoreConfig{
					Debounce: validDebounce,
					DuckLake: DuckLakeConfig{
						Enabled:  true,
						DataPath: "",
					},
				},
			},
			wantErr: "data_path must not be empty",
		},
		{
			name: "enabled with ducklake and data path succeeds",
			cfg: IdentitySignalsConfig{
				Enabled: true,
				Store: StoreConfig{
					Debounce: validDebounce,
					DuckLake: DuckLakeConfig{
						Enabled:  true,
						DataPath: "/var/lib/zitadel/signals",
					},
				},
			},
			wantErr: "",
		},
		{
			name: "S3 backend without bucket fails",
			cfg: IdentitySignalsConfig{
				Enabled: true,
				Store: StoreConfig{
					Debounce: validDebounce,
					DuckLake: DuckLakeConfig{
						Enabled:  true,
						DataPath: "/data",
						Backend:  ArchiveBackendS3,
						S3:       ArchiveS3Config{Bucket: ""},
					},
				},
			},
			wantErr: "s3 bucket must not be empty",
		},
		{
			name: "S3 backend with bucket succeeds",
			cfg: IdentitySignalsConfig{
				Enabled: true,
				Store: StoreConfig{
					Debounce: validDebounce,
					DuckLake: DuckLakeConfig{
						Enabled:  true,
						DataPath: "s3://signals",
						Backend:  ArchiveBackendS3,
						S3:       ArchiveS3Config{Bucket: "my-bucket"},
					},
				},
			},
			wantErr: "",
		},
		{
			name: "zero min_frequency fails",
			cfg: IdentitySignalsConfig{
				Enabled: true,
				Store: StoreConfig{
					Debounce: DebouncerConfig{MinFrequency: 0, MaxBulkSize: 100},
					DuckLake: DuckLakeConfig{Enabled: true, DataPath: "/data"},
				},
			},
			wantErr: "min_frequency must be > 0",
		},
		{
			name: "zero max_bulk_size fails",
			cfg: IdentitySignalsConfig{
				Enabled: true,
				Store: StoreConfig{
					Debounce: DebouncerConfig{MinFrequency: time.Second, MaxBulkSize: 0},
					DuckLake: DuckLakeConfig{Enabled: true, DataPath: "/data"},
				},
			},
			wantErr: "max_bulk_size must be > 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestStreamsConfig_RetentionForStream(t *testing.T) {
	cfg := StreamsConfig{
		Requests: StreamConfig{Retention: 720 * time.Hour},
		Events:   StreamConfig{Retention: 2160 * time.Hour},
	}
	if got := cfg.RetentionForStream(StreamRequests); got != 720*time.Hour {
		t.Errorf("requests retention = %v, want 720h", got)
	}
	if got := cfg.RetentionForStream(StreamEvents); got != 2160*time.Hour {
		t.Errorf("events retention = %v, want 2160h", got)
	}
	if got := cfg.RetentionForStream("unknown"); got != 0 {
		t.Errorf("unknown stream retention = %v, want 0", got)
	}
}

func TestStreamsConfig_EnabledStreams(t *testing.T) {
	tests := []struct {
		name string
		cfg  StreamsConfig
		want int
	}{
		{"both enabled", StreamsConfig{
			Requests: StreamConfig{Enabled: true},
			Events:   StreamConfig{Enabled: true},
		}, 2},
		{"only requests", StreamsConfig{
			Requests: StreamConfig{Enabled: true},
			Events:   StreamConfig{Enabled: false},
		}, 1},
		{"none enabled", StreamsConfig{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.EnabledStreams()
			if len(got) != tt.want {
				t.Errorf("EnabledStreams() len = %d, want %d", len(got), tt.want)
			}
		})
	}
}
