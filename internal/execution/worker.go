package execution

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
)

//go:embed remove_event_executions.sql
var removeEventExecution string

type Worker struct {
	client  *database.DB
	config  WorkerConfig
	queries *ExecutionsQueries
	now     nowFunc
}

// nowFunc makes [time.Now] mockable
type nowFunc func() time.Time

type WorkerConfig struct {
	Workers             uint8
	BulkLimit           uint16
	RequeueEvery        time.Duration
	TransactionDuration time.Duration
	MaxTtl              time.Duration
}

func NewWorker(
	config WorkerConfig,
	client *database.DB,
	queries *ExecutionsQueries,
) *Worker {
	return &Worker{
		config:  config,
		client:  client,
		queries: queries,
		now:     time.Now,
	}
}

func (w *Worker) Start(ctx context.Context) {
	for i := 0; i < int(w.config.Workers); i++ {
		go w.schedule(ctx, i)
	}
}

func (w *Worker) schedule(ctx context.Context, workerID int) {
	t := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			w.log(workerID).Info("scheduler stopped")
			return
		case <-t.C:
			instances, err := w.queryInstances(ctx)
			w.log(workerID).OnError(err).Error("unable to query instances")

			w.triggerInstances(call.WithTimestamp(ctx), instances, workerID)
			t.Reset(w.config.RequeueEvery)
		}
	}
}

func (w *Worker) log(workerID int) *logging.Entry {
	return logging.WithFields("notification worker", workerID)
}

func (w *Worker) queryInstances(ctx context.Context) ([]string, error) {
	return w.queries.ActiveInstancesWithFeatureFlag(ctx), nil
}

func (w *Worker) triggerInstances(ctx context.Context, instances []string, workerID int) {
	for _, instance := range instances {
		instanceCtx := authz.WithInstanceID(ctx, instance)

		err := w.trigger(instanceCtx, workerID)
		w.log(workerID).WithField("instance", instance).OnError(err).Info("trigger failed")
	}
}

func (w *Worker) trigger(ctx context.Context, workerID int) (err error) {
	txCtx := ctx
	if w.config.TransactionDuration > 0 {
		var cancel, cancelTx func()
		txCtx, cancelTx = context.WithCancel(ctx)
		defer cancelTx()
		ctx, cancel = context.WithTimeout(ctx, w.config.TransactionDuration)
		defer cancel()
	}
	tx, err := w.client.BeginTx(txCtx, nil)
	if err != nil {
		return err
	}
	defer func() {
		err = database.CloseTransaction(tx, err)
	}()

	eventExecutions, err := w.queries.searchEventExecutions(txCtx)
	if err != nil {
		return err
	}

	// If there aren't any events or no unlocked event terminate early and start a new run.
	if eventExecutions == nil || len(eventExecutions.EventExecutions) == 0 {
		return nil
	}

	w.log(workerID).
		WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
		WithField("events", len(eventExecutions.EventExecutions)).
		Info("handling event executions")

	for _, event := range eventExecutions.EventExecutions {
		if err := w.removeEventExecution(ctx, tx, event, workerID); err != nil {
			event.WithLogFields(w.log(workerID).OnError(err)).Error("could not handle removal of execution event entry")
			return err
		}

		if err := w.reduceEventExecution(ctx, event); err != nil {
			event.WithLogFields(w.log(workerID).OnError(err)).Error("could not handle event execution")
			return err
		}

		// if the context is canceled, we stop the processing
		if err := ctx.Err(); err != nil {
			event.WithLogFields(w.log(workerID).OnError(err)).Error("timeout handle event executions")
			return nil
		}
	}

	w.log(workerID).
		WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
		WithField("events", len(eventExecutions.EventExecutions)).
		Info("end handling event executions")

	return nil
}

func (w *Worker) reduceEventExecution(ctx context.Context, event *EventExecution) (err error) {
	ctx = ContextWithExecuter(ctx, event.Aggregate)

	// if the notification is too old, we can directly return as it will be removed anyway
	if event.CreatedAt.Add(w.config.MaxTtl).Before(w.now()) {
		return nil
	}

	targets, err := event.Targets()
	if err != nil {
		return err
	}

	_, err = CallTargets(ctx, targets, event.ContextInfo())
	return err
}

func (w *Worker) removeEventExecution(ctx context.Context, tx *sql.Tx, event *EventExecution, workerID int) error {
	_, err := tx.ExecContext(ctx, removeEventExecution,
		event.Aggregate.InstanceID,
		event.Aggregate.ResourceOwner,
		event.Aggregate.Type,
		event.Aggregate.Version,
		event.Aggregate.ID,
		event.Sequence,
		event.EventType,
	)
	event.WithLogFields(w.log(workerID).OnError(err)).Error("could not remove event execution for event")
	return err
}

var _ ContextInfo = &ContextInfoEvent{}

type ContextInfoEvent struct {
	AggregateID   string `json:"aggregateID,omitempty"`
	AggregateType string `json:"aggregateType,omitempty"`
	ResourceOwner string `json:"resourceOwner,omitempty"`
	InstanceID    string `json:"instanceID,omitempty"`
	Version       string `json:"version,omitempty"`
	Sequence      uint64 `json:"sequence,omitempty"`
	EventType     string `json:"event_type,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UserID        string `json:"userID,omitempty"`
	EventPayload  []byte `json:"event_payload,omitempty"`
}

func (c *ContextInfoEvent) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoEvent) SetHTTPResponseBody(resp []byte) error {
	// response is irrelevant and will not be unmarshalled
	return nil
}

func (c *ContextInfoEvent) GetContent() interface{} {
	return c.EventPayload
}
