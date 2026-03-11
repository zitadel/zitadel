package detection

import "github.com/zitadel/zitadel/internal/signals"

type Finding struct {
	Name       string
	Source     string
	Message    string
	Block      bool
	Confidence float64
	// Challenge indicates this finding requires a user challenge (e.g. captcha)
	// rather than an outright block.
	Challenge     bool
	ChallengeType string // e.g. "captcha"
}

type Decision struct {
	Allow    bool
	Findings []Finding
}

// HasChallenge returns true if any finding requires a user challenge.
func (d Decision) HasChallenge() bool {
	for _, f := range d.Findings {
		if f.Challenge {
			return true
		}
	}
	return false
}

// HasBlockingFindings returns true if any finding is a hard block.
func (d Decision) HasBlockingFindings() bool {
	for _, f := range d.Findings {
		if f.Block {
			return true
		}
	}
	return false
}

// ChallengeType returns the type of the first challenge finding, or empty.
func (d Decision) ChallengeType() string {
	for _, f := range d.Findings {
		if f.Challenge {
			return f.ChallengeType
		}
	}
	return ""
}

func (d Decision) BlockingFindings() []Finding {
	findings := make([]Finding, 0, len(d.Findings))
	for _, finding := range d.Findings {
		if finding.Block {
			findings = append(findings, finding)
		}
	}
	return findings
}

func (d Decision) ChallengeFindings() []Finding {
	findings := make([]Finding, 0, len(d.Findings))
	for _, finding := range d.Findings {
		if finding.Challenge {
			findings = append(findings, finding)
		}
	}
	return findings
}

func recordedFindings(findings []Finding) []signals.RecordedFinding {
	if len(findings) == 0 {
		return nil
	}
	recorded := make([]signals.RecordedFinding, len(findings))
	for i, finding := range findings {
		recorded[i] = signals.RecordedFinding{
			Name:          finding.Name,
			Source:        finding.Source,
			Message:       finding.Message,
			Block:         finding.Block,
			Confidence:    finding.Confidence,
			Challenge:     finding.Challenge,
			ChallengeType: finding.ChallengeType,
		}
	}
	return recorded
}

func findingFromRecorded(finding signals.RecordedFinding) Finding {
	return Finding{
		Name:          finding.Name,
		Source:        finding.Source,
		Message:       finding.Message,
		Block:         finding.Block,
		Confidence:    finding.Confidence,
		Challenge:     finding.Challenge,
		ChallengeType: finding.ChallengeType,
	}
}
