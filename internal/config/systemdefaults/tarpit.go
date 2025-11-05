package systemdefaults

import "time"

type TarpitConfig struct {
	// After how many failed attempts, the tarpit should start.
	MinFailedAttempts uint64
	// The seconds that will be added per step.
	StepDuration time.Duration
	// The failed attempts that are needed to increase the tarpit by one step.
	StepSize uint64
	// The maximum duration the tarpit can reach.
	MaxDuration time.Duration
}

func (t *TarpitConfig) Tarpit() func(failedCount uint64) {
	return func(failedCount uint64) {
		time.Sleep(t.duration(failedCount))
	}
}

func (t *TarpitConfig) duration(failedCount uint64) time.Duration {
	if failedCount < t.MinFailedAttempts {
		return 0
	}
	// calculate the step we are at
	// every StepSize failed attempts increase the step by one
	step := (failedCount - t.MinFailedAttempts) / t.StepSize
	duration := time.Duration(step) * t.StepDuration
	if duration < t.MaxDuration {
		return duration
	}
	return t.MaxDuration
}
