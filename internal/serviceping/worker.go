package serviceping

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/riverqueue/river"
	"github.com/robfig/cron/v3"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/queue"
	"github.com/zitadel/zitadel/internal/v2/system"
	"github.com/zitadel/zitadel/internal/zerrors"
	analytics "github.com/zitadel/zitadel/pkg/grpc/analytics/v2beta"
)

const (
	QueueName   = "service_ping_report"
	minInterval = 30 * time.Minute
)

var (
	ErrInvalidReportType = errors.New("invalid report type")

	_ river.Worker[*ServicePingReport] = (*Worker)(nil)
)

type Worker struct {
	river.WorkerDefaults[*ServicePingReport]

	reportClient analytics.TelemetryServiceClient
	db           Queries
	queue        Queue

	config   *Config
	systemID string
	version  string
}

type Queries interface {
	SearchInstances(ctx context.Context, queries *query.InstanceSearchQueries) (*query.Instances, error)
	ListResourceCounts(ctx context.Context, lastID int, size int) ([]query.ResourceCount, error)
}

type Queue interface {
	Insert(ctx context.Context, args river.JobArgs, opts ...queue.InsertOpt) error
}

// Register implements the [queue.Worker] interface.
func (w *Worker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker[*ServicePingReport](workers, w)
	queues[QueueName] = river.QueueConfig{
		MaxWorkers: 1, // for now, we only use a single worker to prevent too much side effects on other queues
	}
}

// Work implements the [river.Worker] interface.
func (w *Worker) Work(ctx context.Context, job *river.Job[*ServicePingReport]) (err error) {
	defer func() {
		err = w.handleClientError(err)
	}()
	switch job.Args.ReportType {
	case ReportTypeBaseInformation:
		reportID, err := w.reportBaseInformation(ctx)
		if err != nil {
			return err
		}
		return w.createReportJobs(ctx, reportID)
	case ReportTypeResourceCounts:
		return w.reportResourceCounts(ctx, job.Args.ReportID)
	default:
		logging.WithFields("reportType", job.Args.ReportType, "reportID", job.Args.ReportID).
			Error("unknown job type")
		return river.JobCancel(ErrInvalidReportType)
	}
}

func (w *Worker) reportBaseInformation(ctx context.Context) (string, error) {
	instances, err := w.db.SearchInstances(ctx, &query.InstanceSearchQueries{})
	if err != nil {
		return "", err
	}
	instanceInformation := instanceInformationToPb(instances)
	resp, err := w.reportClient.ReportBaseInformation(ctx, &analytics.ReportBaseInformationRequest{
		SystemId:  w.systemID,
		Version:   w.version,
		Instances: instanceInformation,
	})
	if err != nil {
		return "", err
	}
	return resp.GetReportId(), nil
}

func (w *Worker) reportResourceCounts(ctx context.Context, reportID string) error {
	lastID := 0
	// iterate over the resource counts until there are no more counts to report
	// or the context gets cancelled
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			counts, err := w.db.ListResourceCounts(ctx, lastID, w.config.Telemetry.ResourceCount.BulkSize)
			if err != nil {
				return err
			}
			// if there are no counts, we can stop the loop
			if len(counts) == 0 {
				return nil
			}
			request := &analytics.ReportResourceCountsRequest{
				SystemId:       w.systemID,
				ResourceCounts: resourceCountsToPb(counts),
			}
			if reportID != "" {
				request.ReportId = gu.Ptr(reportID)
			}
			resp, err := w.reportClient.ReportResourceCounts(ctx, request)
			if err != nil {
				return err
			}
			// in case the resource counts returned by the database are less than the bulk size,
			// we can assume that we have reached the end of the resource counts and can stop the loop
			if len(counts) < w.config.Telemetry.ResourceCount.BulkSize {
				return nil
			}
			// update the lastID for the next iteration
			lastID = counts[len(counts)-1].ID
			// In case we get a report ID back from the server (it could be the first call of the report),
			// we update it to use it for the next batch.
			if resp.GetReportId() != "" && resp.GetReportId() != reportID {
				reportID = resp.GetReportId()
			}
		}
	}
}

func (w *Worker) handleClientError(err error) error {
	telemetryError := new(TelemetryError)
	if !errors.As(err, &telemetryError) {
		// If the error is not a TelemetryError, we can assume that it is a transient error
		// and can be retried by the queue.
		return err
	}
	switch telemetryError.StatusCode {
	case http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusNotImplemented,
		http.StatusConflict,
		http.StatusPreconditionFailed:
		// In case of these errors, we can assume that a retry does not make sense,
		// so we can cancel the job.
		return river.JobCancel(err)
	default:
		// As of now we assume that all other errors are transient and can be retried.
		// So we just return the error, which will be handled by the queue as a failed attempt.
		return err
	}
}

func (w *Worker) createReportJobs(ctx context.Context, reportID string) error {
	errs := make([]error, 0)
	if w.config.Telemetry.ResourceCount.Enabled {
		err := w.addReportJob(ctx, reportID, ReportTypeResourceCounts)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (w *Worker) addReportJob(ctx context.Context, reportID string, reportType ReportType) error {
	job := &ServicePingReport{
		ReportID:   reportID,
		ReportType: reportType,
	}
	return w.queue.Insert(ctx, job,
		queue.WithQueueName(QueueName),
		queue.WithMaxAttempts(w.config.MaxAttempts),
	)
}

type systemIDReducer struct {
	id string
}

func (s *systemIDReducer) Reduce() error {
	return nil
}

func (s *systemIDReducer) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		if idEvent, ok := event.(*system.IDGeneratedEvent); ok {
			s.id = idEvent.ID
		}
	}
}

func (s *systemIDReducer) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(system.AggregateType).
		EventTypes(system.IDGeneratedType).
		Builder()
}

func Register(
	ctx context.Context,
	q *queue.Queue,
	queries *query.Queries,
	eventstoreClient *eventstore.Eventstore,
	config *Config,
) error {
	if !config.Enabled {
		return nil
	}
	systemID := new(systemIDReducer)
	err := eventstoreClient.FilterToQueryReducer(ctx, systemID)
	if err != nil {
		return err
	}
	q.AddWorkers(&Worker{
		reportClient: NewClient(config),
		db:           queries,
		queue:        q,
		config:       config,
		systemID:     systemID.id,
		version:      build.Version(),
	})
	return nil
}

func Start(config *Config, q *queue.Queue) error {
	if !config.Enabled {
		return nil
	}
	schedule, err := parseAndValidateSchedule(config.Interval)
	if err != nil {
		return err
	}
	q.AddPeriodicJob(
		schedule,
		&ServicePingReport{},
		queue.WithQueueName(QueueName),
		queue.WithMaxAttempts(config.MaxAttempts),
	)
	return nil
}

func parseAndValidateSchedule(interval string) (cron.Schedule, error) {
	if interval == "@daily" {
		interval = randomizeDaily()
	}
	schedule, err := cron.ParseStandard(interval)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "SERV-NJqiof", "invalid interval")
	}
	var intervalDuration time.Duration
	switch s := schedule.(type) {
	case *cron.SpecSchedule:
		// For cron.SpecSchedule, we need to calculate the interval duration
		// by getting the next time and subtracting it from the time after that.
		// This is because the schedule could be a specific time, that is less than 30 minutes away,
		// but still run only once a day and therefore is valid.
		next := s.Next(time.Now())
		nextAfter := s.Next(next)
		intervalDuration = nextAfter.Sub(next)
	case cron.ConstantDelaySchedule:
		intervalDuration = s.Delay
	}
	if intervalDuration < minInterval {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "SERV-FJ12", "interval must be at least %s", minInterval)
	}
	logging.WithFields("interval", interval).Info("scheduling service ping")
	return schedule, nil
}

// randomizeDaily generates a random time for the daily cron job
// to prevent all systems from sending the report at the same time.
func randomizeDaily() string {
	minute := rand.Intn(60)
	hour := rand.Intn(24)
	return fmt.Sprintf("%d %d * * *", minute, hour)
}
