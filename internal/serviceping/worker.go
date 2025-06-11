package serviceping

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/muhlemmer/gu"
	"github.com/riverqueue/river"
	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/queue"
	"github.com/zitadel/zitadel/internal/v2/system"
	analytics "github.com/zitadel/zitadel/pkg/grpc/analytics/v2beta"
)

const (
	QueueName = "service_ping_report"
)

type ReportType uint

type Worker struct {
	river.WorkerDefaults[*ServicePingReport]

	reportClient analytics.TelemetryServiceClient
	db           *query.Queries
	queue        *queue.Queue

	config   *Config
	systemID string

	//now    nowFunc
}

// Work implements [river.Worker].
func (w *Worker) Work(ctx context.Context, job *river.Job[*ServicePingReport]) (err error) {
	defer func() {
		err = w.handleClientError(err)
	}()
	switch job.Args.TestState {
	case 0:
		reportID, err := w.reportBaseInformation(ctx)
		if err != nil {
			return err
		}
		return w.createReportJobs(ctx, reportID)
	case 1:
		return w.reportResourceCounts(ctx, job.Args.ReportID)
	}
	return nil
}

func (w *Worker) reportBaseInformation(ctx context.Context) (string, error) {
	instances, err := w.db.SearchInstances(ctx, &query.InstanceSearchQueries{})
	if err != nil {
		return "", err
	}
	instanceInformation := instanceInformationToPb(instances)

	resp, err := w.reportClient.ReportBaseInformation(ctx, &analytics.ReportBaseInformationRequest{
		SystemId:  w.systemID,
		Version:   build.Version(),
		Instances: instanceInformation,
	})
	if err != nil {
		return "", err
	}
	return resp.ReportId, nil
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
			resourceCounts := resourceCountsToPb(counts)
			resp, err := w.reportClient.ReportResourceCounts(ctx, &analytics.ReportResourceCountsRequest{
				SystemId:       w.systemID,
				ReportId:       gu.Ptr(reportID),
				ResourceCounts: resourceCounts,
			})
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
	switch status.Code(err) {
	case codes.OK:
		return nil
	case codes.InvalidArgument,
		codes.NotFound,
		codes.Unimplemented,
		codes.AlreadyExists,
		codes.FailedPrecondition:
		// In case of these errors, we can assume that a retry does not make sense,
		// so we can cancel the job.
		return river.JobCancel(err)
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Unavailable:
		// These errors are typically transient and can be retried.
		return err
	default:
		return err
	}
}

type ServicePingReport struct {
	ReportID  string
	TestState int
}

func (r *ServicePingReport) Kind() string {
	return "service_ping_report"
}

func NewWorker(
	reportClient analytics.TelemetryServiceClient,
	queries *query.Queries,
	q *queue.Queue,
	config *Config,
	systemID string,

) *Worker {
	return &Worker{
		reportClient: reportClient,
		db:           queries,
		queue:        q,
		config:       config,
		systemID:     systemID,
	}
}

var _ river.Worker[*ServicePingReport] = (*Worker)(nil)

func (w *Worker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker[*ServicePingReport](workers, w)
	queues[QueueName] = river.QueueConfig{
		//MaxWorkers: int(w.config.Workers),
		MaxWorkers: 1,
	}
}

func (w *Worker) createReportJobs(ctx context.Context, reportID string) error {
	if w.config.Telemetry.ResourceCount.Enabled {
		err := w.addReportJob(ctx, reportID, 1)
		logging.WithFields("reportID", reportID, "telemetry", "resourceCount").
			OnError(err).Error("failed to add report job")
	}
	return nil
}

func (w *Worker) addReportJob(ctx context.Context, reportID string, state int) error {
	job := &ServicePingReport{
		ReportID:  reportID,
		TestState: state,
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

func StartWorker(
	ctx context.Context,
	q *queue.Queue,
	queries *query.Queries,
	eventstoreClient *eventstore.Eventstore,
	config *Config,
) error {
	systemID := new(systemIDReducer)
	err := eventstoreClient.FilterToQueryReducer(ctx, systemID)
	if err != nil {
		return err
	}
	q.AddWorkers(NewWorker(
		&mockReportClient{},
		queries,
		q,
		config,
		systemID.id,
	))
	return nil
}

type mockReportClient struct{}

func (m mockReportClient) ReportBaseInformation(ctx context.Context, in *analytics.ReportBaseInformationRequest, opts ...grpc.CallOption) (*analytics.ReportBaseInformationResponse, error) {
	fmt.Println(proto.MarshalTextString(in))
	reportID, err := id.SonyFlakeGenerator().Next()
	if err != nil {
		return nil, err
	}
	return &analytics.ReportBaseInformationResponse{
		ReportId: reportID,
	}, nil
}

func (m mockReportClient) ReportResourceCounts(ctx context.Context, in *analytics.ReportResourceCountsRequest, opts ...grpc.CallOption) (*analytics.ReportResourceCountsResponse, error) {
	fmt.Println(proto.MarshalTextString(in))
	return &analytics.ReportResourceCountsResponse{
		ReportId: in.GetReportId(),
	}, nil
}
