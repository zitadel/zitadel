package grpc

type Config struct {
	ServerPort    string
	GatewayPort   string
	CustomHeaders []string
}

func (c Config) ToServerConfig() ServerConfig {
	return ServerConfig{
		Port: c.ServerPort,
	}
}

func (c Config) ToGatewayConfig() GatewayConfig {
	return GatewayConfig{
		Port:          c.GatewayPort,
		GRPCEndpoint:  c.ServerPort,
		CustomHeaders: c.CustomHeaders,
	}
}

type ServerConfig struct {
	Port string
}

type GatewayConfig struct {
	Port          string
	GRPCEndpoint  string
	CustomHeaders []string
}
