package serviceping

type Config struct {
	Enabled     bool
	Endpoint    string
	Interval    string
	MaxAttempts uint8
	Telemetry   TelemetryConfig
}

type TelemetryConfig struct {
	ResourceCount ResourceCount
}

type ResourceCount struct {
	Enabled  bool
	BulkSize int
}
