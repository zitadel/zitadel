package serviceping

type Config struct {
	Enabled   bool
	Endpoint  string
	Interval  string
	Telemetry TelemetryConfig
}

type TelemetryConfig struct {
	ResourceCount struct {
		Enabled  bool
		BulkSize uint
	}
}
