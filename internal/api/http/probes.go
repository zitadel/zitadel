package http

const (
	Healthz    = "/healthz"
	Readiness  = "/ready"
	Validation = "/validate"
)

var (
	Probes = []string{Healthz, Readiness, Validation}
)
