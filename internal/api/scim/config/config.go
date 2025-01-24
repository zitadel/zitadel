package config

type Config struct {
	DocumentationUrl      string
	AuthenticationSchemes []*ServiceProviderConfigAuthenticationScheme
	EmailVerified         bool
	PhoneVerified         bool
	MaxRequestBodySize    int64
	Bulk                  BulkConfig
}

type BulkConfig struct {
	MaxOperationsCount int
}

type ServiceProviderConfigAuthenticationScheme struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	SpecUri          string `json:"specUri"`
	DocumentationUri string `json:"documentationUri"`
	Type             string `json:"type"`
	Primary          bool   `json:"primary"`
}
