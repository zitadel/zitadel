package logstore

type Configs struct {
	Access    *Config
	Execution *Config
}

type Config struct {
	Database *EmitterConfig
	Stdout   *EmitterConfig
}
