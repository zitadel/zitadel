package logstore

type Configs struct {
	Access    *Config
	Execution *Config
}

type Config struct {
	Stdout *StdConfig
}

type StdConfig struct {
	Enabled bool
}
