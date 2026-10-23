package queue

import (
	"context"
	"testing"

	"github.com/riverqueue/river"
)

// testJobArgs is a minimal river.JobArgs implementation for tests.
type testJobArgs struct{}

func (testJobArgs) Kind() string { return "test" }

// A *Queue zero value has no river client; tests use the nil receiver form
// below to stay independent of a database.
func TestQueueInsertNilReceiver(t *testing.T) {
	// Regression test for #12887: cmd/mirror passes a nil *Queue into
	// notification.Register, so the back-channel logout reducer (and any
	// other notification reducer) calls Insert on a nil receiver while
	// projections are recomputed. This used to panic with a nil pointer
	// dereference at queue.go:118, which killed `zitadel mirror` with
	// exit code 2 before the verify step. Insert must instead be a no-op,
	// like the other methods on *Queue that tolerate a nil receiver
	// (ShouldStart, AddWorkers, AddPeriodicJob).
	var q *Queue
	if err := q.Insert(context.Background(), testJobArgs{}, WithQueueName("test")); err != nil {
		t.Errorf("Insert on nil receiver must not fail: %v", err)
	}
}

func TestQueueInsertManyFastTxNilReceiver(t *testing.T) {
	// Same nil-receiver contract for the transactional bulk insert used by
	// the eventstore push path.
	var q *Queue
	if err := q.InsertManyFastTx(context.Background(), nil, []river.JobArgs{testJobArgs{}}); err != nil {
		t.Errorf("InsertManyFastTx on nil receiver must not fail: %v", err)
	}
}
