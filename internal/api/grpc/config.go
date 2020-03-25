package grpc

type Config struct {
	ServerPort    string
	GatewayPort   string
	SearchLimit   int
	CustomHeaders []string
}

func (c *Config) ToServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:        c.ServerPort,
		SearchLimit: c.SearchLimit,
	}
}

func (c *Config) ToGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Port:          c.GatewayPort,
		GRPCEndpoint:  c.ServerPort,
		CustomHeaders: c.CustomHeaders,
	}
}

type ServerConfig struct {
	Port        string
	SearchLimit int
}

type GatewayConfig struct {
	Port          string
	GRPCEndpoint  string
	CustomHeaders []string
}
