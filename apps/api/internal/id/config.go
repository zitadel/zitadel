package id

const (
	DefaultWebhookPath = "http://metadata.google.internal/computeMetadata/v1/instance/id"
)

type Config struct {
	// Configuration for the identification of machines.
	Identification Identification
}

type Identification struct {
	// Configuration for using private IP to identify a machine.
	PrivateIp PrivateIp
	// Configuration for using hostname to identify a machine.
	Hostname Hostname
	// Configuration for using a webhook to identify a machine.
	Webhook Webhook
}

type PrivateIp struct {
	// Try to use private IP when identifying the machine uniquely
	Enabled bool
}

type Hostname struct {
	// Try to use hostname when identifying the machine uniquely
	Enabled bool
}

type Webhook struct {
	// Try to use webhook when identifying the machine uniquely
	Enabled bool
	// The URL of the metadata endpoint to query
	Url string
	// (Optional) A JSONPath expression for the data to extract from the response from the metadata endpoint
	JPath *string
	// (Optional) Headers to pass in the metadata request
	Headers *map[string]string
}

func Configure(config *Config) {
	if config != nil {
		GeneratorConfig = config
	}
}
