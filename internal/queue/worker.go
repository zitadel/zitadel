package queue

// import (
// 	"context"
// 	"time"

// 	"github.com/riverqueue/river"
// 	"github.com/riverqueue/river/rivertype"
// )

// type testArgs struct {
// 	Arg1 string
// }

// func (testArgs) Kind() string {
// 	return "test"
// }

// type testWorker struct {
// 	river.WorkerDefaults[testArgs]
// }

// // Middleware implements [river.Worker].
// // Subtle: this method shadows the method (WorkerDefaults).Middleware of testWorker.WorkerDefaults.
// func (t *testWorker) Middleware(job *river.Job[testArgs]) []rivertype.WorkerMiddleware {
// 	panic("unimplemented")
// }

// // NextRetry implements [river.Worker].
// // Subtle: this method shadows the method (WorkerDefaults).NextRetry of testWorker.WorkerDefaults.
// func (t *testWorker) NextRetry(job *river.Job[testArgs]) time.Time {
// 	panic("unimplemented")
// }

// // Timeout implements [river.Worker].
// // Subtle: this method shadows the method (WorkerDefaults).Timeout of testWorker.WorkerDefaults.
// func (t *testWorker) Timeout(job *river.Job[testArgs]) time.Duration {
// 	panic("unimplemented")
// }

// // Work implements [river.Worker].
// func (t *testWorker) Work(ctx context.Context, job *river.Job[testArgs]) error {
// 	panic("unimplemented")
// }

// func (t *testWorker) Register(workers *river.Workers) {
// 	river.AddWorker(workers, t)
// }

// var _ river.Worker[testArgs] = (*testWorker)(nil)

// func bla() {
// 	q := &Queue{}
// 	q.AddWorkers(&testWorker{})
// }
