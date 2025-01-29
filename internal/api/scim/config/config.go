package config

type Config struct {
	EmailVerified      bool
	PhoneVerified      bool
	MaxRequestBodySize int64
	Bulk               BulkConfig
}

type BulkConfig struct {
	MaxOperationsCount int
}
