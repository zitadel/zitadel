package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/internal/query"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
)

type Worker struct {
	river.WorkerDefaults[*exec_repo.Request]

	config WorkerConfig
	now    nowFunc
}

// Timeout implements the Timeout-function of [river.Worker].
// Maximum time a job can run before the context gets cancelled.
// The time can be shorter than the sum of target timeouts, this is expected behavior to not block the request indefinitely.
func (w *Worker) Timeout(*river.Job[*exec_repo.Request]) time.Duration {
	return w.config.TransactionDuration
}

// Work implements [river.Worker].
func (w *Worker) Work(ctx context.Context, job *river.Job[*exec_repo.Request]) error {
	ctx = ContextWithExecuter(ctx, job.Args.Aggregate)

	// if the event is too old, we can directly return as it will be removed anyway
	if job.CreatedAt.Add(w.config.MaxTtl).Before(w.now()) {
		return river.JobCancel(errors.New("event is too old"))
	}

	targets, err := TargetsFromRequest(job.Args)
	if err != nil {
		// If we are not able to get the targets from the request, we can cancel the job, as we have nothing to call
		return river.JobCancel(fmt.Errorf("unable to unmarshal targets because %w", err))
	}

	_, err = CallTargets(ctx, targets, exec_repo.ContextInfoFromRequest(job.Args))
	if err != nil {
		// If there is an error returned from the targets, it means that the execution was interrupted
		return river.JobCancel(fmt.Errorf("interruption during call of targets because %w", err))
	}
	return nil
}

// nowFunc makes [time.Now] mockable
type nowFunc func() time.Time

type WorkerConfig struct {
	Workers             uint8
	TransactionDuration time.Duration
	MaxTtl              time.Duration
}

func NewWorker(
	config WorkerConfig,
) *Worker {
	return &Worker{
		config: config,
		now:    time.Now,
	}
}

var _ river.Worker[*exec_repo.Request] = (*Worker)(nil)

func (w *Worker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker(workers, w)
	queues[exec_repo.QueueName] = river.QueueConfig{
		MaxWorkers: int(w.config.Workers),
	}
}

func TargetsFromRequest(e *exec_repo.Request) ([]Target, error) {
	var execTargets []*query.ExecutionTarget
	if err := json.Unmarshal(e.TargetsData, &execTargets); err != nil {
		return nil, err
	}
	targets := make([]Target, len(execTargets))
	for i, target := range execTargets {
		targets[i] = target
	}
	return targets, nil
}
