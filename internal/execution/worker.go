package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/internal/api/oidc/sign"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/denylist"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
)

type Worker struct {
	river.WorkerDefaults[*exec_repo.Request]

	config     WorkerConfig
	httpClient *http.Client
	now        NowFunc

	targetEncAlg     crypto.EncryptionAlgorithm
	activeSigningKey GetActiveSigningWebKey
}

// Timeout implements the Timeout-function of [river.Worker].
// Maximum time a job can run before the context gets canceled.
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

	info := exec_repo.ContextInfoFromRequest(job.Args)
	signerOnce := sign.GetSignerOnce(w.activeSigningKey)
	encrypters := &sync.Map{}
	var deliveryErrors error
	for _, target := range targets {
		// The queue already provides asynchronous execution; acknowledge only after delivery finishes.
		if target.TargetType == target_domain.TargetTypeAsync {
			target.TargetType = target_domain.TargetTypeWebhook
		}
		_, err := CallTarget(ctx, target, info, w.targetEncAlg, signerOnce, encrypters, w.httpClient)
		if err == nil {
			continue
		}
		deliveryErrors = errors.Join(deliveryErrors, fmt.Errorf("target %s: %w", target.GetTargetID(), err))
		if target.IsInterruptOnError() {
			break
		}
	}
	// River retries the whole target list; successful receivers may see the event again.
	return deliveryErrors
}

// NowFunc makes [time.Now] mockable
type NowFunc func() time.Time

type WorkerConfig struct {
	Workers             uint8
	TransactionDuration time.Duration
	MaxTtl              time.Duration
	DenyList            []denylist.AddressChecker
}

func NewWorker(
	config WorkerConfig,
	targetEncAlg crypto.EncryptionAlgorithm,
	activeSigningKey GetActiveSigningWebKey,
	now NowFunc,
	httpClient *http.Client,
) *Worker {
	return &Worker{
		config:           config,
		httpClient:       httpClient,
		now:              now,
		targetEncAlg:     targetEncAlg,
		activeSigningKey: activeSigningKey,
	}
}

var _ river.Worker[*exec_repo.Request] = (*Worker)(nil)

func (w *Worker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker(workers, w)
	queues[exec_repo.QueueName] = river.QueueConfig{
		MaxWorkers: int(w.config.Workers),
	}
}

func TargetsFromRequest(e *exec_repo.Request) ([]target_domain.Target, error) {
	var targets []target_domain.Target
	if err := json.Unmarshal(e.TargetsData, &targets); err != nil {
		return nil, err
	}
	return targets, nil
}
