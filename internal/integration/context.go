package integration

import (
	"context"
	"time"
)

// WaitForAndTickWithMaxDuration determine a duration and interval for EventuallyWithT-tests from context timeout and desired max duration
func WaitForAndTickWithMaxDuration(ctx context.Context, max time.Duration) (time.Duration, time.Duration) {
	// interval which is used to retry the test
	tick := time.Millisecond * 100
	// tolerance which is used to stop the test for the timeout
	tolerance := tick * 5
	// default of the WaitFor is always a defined duration, shortened if the context would time out before
	waitFor := max

	if ctxDeadline, ok := ctx.Deadline(); ok {
		// if the context has a deadline, set the WaitFor to the shorter duration
		if until := time.Until(ctxDeadline); until < waitFor {
			// ignore durations which are smaller than the tolerance
			if until < tolerance {
				waitFor = 0
			} else {
				// always let the test stop with tolerance before the context is in timeout
				waitFor = until - tolerance
			}
		}
	}
	return waitFor, tick
}
