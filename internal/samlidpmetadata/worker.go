package samlidpmetadata

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/robfig/cron/v3"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/queue"
)

const (
	// kind is the job kind used to schedule the periodic metadata refresh.
	kind = "saml_idp_metadata_refresh"
	// queueName is the queue the refresh jobs are processed on.
	queueName = "saml_idp_metadata_refresh"
	// batchSize limits the number of IdPs refreshed per run.
	batchSize = 100
)

// Config configures the periodic SAML IdP metadata refresh.
type Config struct {
	// Interval is the cron schedule on which the SAML IdP metadata is refreshed.
	// An empty interval disables the refresh.
	Interval string
}

var _ river.Worker[*SAMLIDPMetadataRefresh] = (*Worker)(nil)

// SAMLIDPMetadataRefresh is a periodic job which re-fetches the metadata of SAML
// identity providers configured with a metadata URL.
type SAMLIDPMetadataRefresh struct{}

func (SAMLIDPMetadataRefresh) Kind() string {
	return kind
}

type Queries interface {
	ListSAMLIDPMetadataURLs(ctx context.Context, afterInstanceID, afterID string, limit uint32) ([]query.SAMLIDPMetadataURL, error)
}

type Commands interface {
	RefreshInstanceSAMLProviderMetadata(ctx context.Context, id string) error
	RefreshOrgSAMLProviderMetadata(ctx context.Context, resourceOwner, id string) error
}

type Worker struct {
	river.WorkerDefaults[*SAMLIDPMetadataRefresh]

	queries  Queries
	commands Commands
}

func NewWorker(queries Queries, commands Commands) *Worker {
	return &Worker{
		queries:  queries,
		commands: commands,
	}
}

// Register implements the [queue.Worker] interface.
func (w *Worker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker(workers, w)
	queues[queueName] = river.QueueConfig{MaxWorkers: 1}
}

func (w *Worker) Work(ctx context.Context, job *river.Job[*SAMLIDPMetadataRefresh]) error {
	var afterInstanceID, afterID string
	for {
		idps, err := w.queries.ListSAMLIDPMetadataURLs(ctx, afterInstanceID, afterID, batchSize)
		if err != nil {
			return err
		}
		for _, idp := range idps {
			if err := w.refresh(ctx, idp); err != nil {
				logging.WithError(err).WithField("instanceID", idp.InstanceID).WithField("idpID", idp.ID).
					Warn("unable to refresh saml idp metadata")
			}
		}
		if len(idps) < batchSize {
			return nil
		}
		afterInstanceID = idps[len(idps)-1].InstanceID
		afterID = idps[len(idps)-1].ID
	}
}

func (w *Worker) refresh(ctx context.Context, idp query.SAMLIDPMetadataURL) error {
	ctx = authz.WithInstanceID(ctx, idp.InstanceID)

	switch idp.OwnerType {
	case domain.IdentityProviderTypeSystem:
		return w.commands.RefreshInstanceSAMLProviderMetadata(ctx, idp.ID)
	case domain.IdentityProviderTypeOrg:
		return w.commands.RefreshOrgSAMLProviderMetadata(ctx, idp.ResourceOwner, idp.ID)
	default:
		return fmt.Errorf("unsupported identity provider owner type %d", idp.OwnerType)
	}
}

// Register adds the worker to the queue. It must be called before starting the queue.
func Register(ctx context.Context, q *queue.Queue, queries Queries, commands Commands, config *Config) {
	if config == nil || config.Interval == "" {
		return
	}
	q.ShouldStart()
	q.AddWorkers(ctx, NewWorker(queries, commands))
}

// Start schedules the periodic metadata refresh. It must be called after the queue started.
func Start(ctx context.Context, q *queue.Queue, config *Config) error {
	if config == nil || config.Interval == "" {
		return nil
	}
	schedule, err := cron.ParseStandard(config.Interval)
	if err != nil {
		return err
	}
	q.AddPeriodicJob(ctx, schedule, &SAMLIDPMetadataRefresh{}, queue.WithQueueName(queueName))
	return nil
}
