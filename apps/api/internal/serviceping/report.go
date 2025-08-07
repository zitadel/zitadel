package serviceping

type ReportType uint

const (
	ReportTypeBaseInformation ReportType = iota
	ReportTypeResourceCounts
)

type ServicePingReport struct {
	ReportID   string
	ReportType ReportType
}

func (r *ServicePingReport) Kind() string {
	return "service_ping_report"
}
