package risk

import "time"

type Outcome string

type Finding struct {
	Name       string
	Source     string
	Message    string
	Block      bool
	Confidence float64
}

type Decision struct {
	Allow    bool
	Findings []Finding
}

type Signal struct {
	InstanceID    string
	UserID        string
	SessionID     string
	FingerprintID string
	Operation     string
	Outcome       Outcome
	Timestamp     time.Time
	IP            string
	UserAgent     string
}

type RecordedSignal struct {
	Signal
	Findings []Finding
}

type Snapshot struct {
	UserSignals    []RecordedSignal
	SessionSignals []RecordedSignal
}

const (
	OutcomeSuccess Outcome = "success"
	OutcomeFailure Outcome = "failure"
	OutcomeBlocked Outcome = "blocked"
)

func (d Decision) BlockingFindings() []Finding {
	findings := make([]Finding, 0, len(d.Findings))
	for _, finding := range d.Findings {
		if finding.Block {
			findings = append(findings, finding)
		}
	}
	return findings
}
