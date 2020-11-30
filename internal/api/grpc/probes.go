package grpc

const (
	Healthz                 = "/Healthz"
	Readiness               = "/Ready"
	Validation              = "/Validate"
	HealthZMethodAdmin      = ""
	HealthZMethodManagement = ""
	HealthZMethodAuth       = ""
)

var (
	Probes     = []string{Healthz, Readiness, Validation}
	GRPCProbes = []string{HealthZMethodAdmin, HealthZMethodManagement, HealthZMethodAuth}
)
