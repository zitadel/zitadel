package grpc

const (
	Healthz    = "/Healthz"
	Readiness  = "/Ready"
	Validation = "/Validate"
)

var (
	Probes = []string{Healthz, Readiness, Validation}
)
