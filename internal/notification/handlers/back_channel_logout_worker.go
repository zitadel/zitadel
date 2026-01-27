package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/riverqueue/river"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/api/oidc/sign"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification/backchannel"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/queue"
)

type BackChannelLogoutWorker struct {
	river.WorkerDefaults[*backchannel.LogoutRequest]

	commands      Commands
	queries       *NotificationQueries
	eventstore    *eventstore.Eventstore
	queue         Queue
	channels      types.ChannelChains
	config        BackChannelLogoutWorkerConfig
	tokenLifetime time.Duration
	now           nowFunc
	idGenerator   id.Generator
}

// Timeout implements the Timeout-function of [river.Worker].
// Maximum time a job can run before the context gets canceled.
func (w *BackChannelLogoutWorker) Timeout(*river.Job[*backchannel.LogoutRequest]) time.Duration {
	return w.config.TransactionDuration
}

// Work implements [river.Worker].
func (w *BackChannelLogoutWorker) Work(ctx context.Context, job *river.Job[*backchannel.LogoutRequest]) error {
	ctx = ContextWithNotifier(ctx, job.Args.Aggregate)

	ctx, err := enrichCtx(ctx, job.Args.TriggeredAtOrigin)
	if err != nil {
		return channels.NewCancelError(err)
	}

	// if the notification is too old, we can directly cancel
	if job.CreatedAt.Add(w.config.MaxTtl).Before(w.now()) {
		return river.JobCancel(errors.New("back channel logout notification is too old"))
	}

	if job.Args.OIDCSessionID == "" {
		return w.createNotificationJobs(ctx, job.Args)
	}
	return w.sendLogoutRequest(ctx, job.Args)
}

func (w *BackChannelLogoutWorker) createNotificationJobs(ctx context.Context, request *backchannel.LogoutRequest) error {
	sessions := &backChannelLogoutSession{sessionID: request.SessionID}
	err := w.eventstore.FilterToQueryReducer(ctx, sessions)
	if err != nil {
		return err
	}

	for _, oidcSession := range sessions.sessions {
		tokenID, err := w.idGenerator.Next()
		if err != nil {
			return err
		}
		logoutRequest := &backchannel.LogoutRequest{
			Aggregate:            request.Aggregate,
			SessionID:            request.SessionID,
			TriggeredAtOrigin:    request.TriggeredAtOrigin,
			TriggeringEventType:  request.TriggeringEventType,
			TokenID:              tokenID,
			UserID:               oidcSession.UserID,
			OIDCSessionID:        oidcSession.OIDCSessionID,
			ClientID:             oidcSession.ClientID,
			BackChannelLogoutURI: oidcSession.BackChannelLogoutURI,
		}
		err = w.queue.Insert(ctx, logoutRequest,
			queue.WithQueueName(backchannel.QueueName),
			queue.WithMaxAttempts(w.config.MaxAttempts))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *BackChannelLogoutWorker) sendLogoutRequest(ctx context.Context, request *backchannel.LogoutRequest) error {
	getSigner := sign.GetSignerOnce(w.queries.GetActiveSigningWebKey)
	token, err := w.logoutToken(ctx, request, getSigner)
	if err != nil {
		return err
	}
	if err = types.SendSecurityTokenEvent(ctx, set.Config{CallURL: request.BackChannelLogoutURI}, w.channels, &LogoutTokenMessage{LogoutToken: token}, request.TriggeringEventType).WithoutTemplate(); err != nil {
		return err
	}
	return w.commands.BackChannelLogoutSent(ctx, request.SessionID, request.OIDCSessionID, request.Aggregate.InstanceID)
}

func (w *BackChannelLogoutWorker) logoutToken(ctx context.Context, request *backchannel.LogoutRequest, getSigner sign.SignerFunc) (string, error) {
	token := oidc.NewLogoutTokenClaims(
		request.TriggeredAtOrigin,
		request.UserID,
		oidc.Audience{request.ClientID},
		w.now().Add(w.tokenLifetime),
		request.TokenID,
		request.SessionID,
		time.Second,
	)
	signer, _, err := getSigner(ctx)
	if err != nil {
		return "", err
	}
	return crypto.Sign(token, signer)
}

func NewBackChannelLogoutWorker(
	commands Commands,
	queries *NotificationQueries,
	eventstore *eventstore.Eventstore,
	queue Queue,
	channels types.ChannelChains,
	config BackChannelLogoutWorkerConfig,
	tokenLifetime time.Duration,
	idGenerator id.Generator,
) *BackChannelLogoutWorker {
	return &BackChannelLogoutWorker{
		commands:      commands,
		queries:       queries,
		eventstore:    eventstore,
		queue:         queue,
		channels:      channels,
		config:        config,
		tokenLifetime: tokenLifetime,
		now:           time.Now,
		idGenerator:   idGenerator,
	}
}

var _ river.Worker[*backchannel.LogoutRequest] = (*BackChannelLogoutWorker)(nil)

func (w *BackChannelLogoutWorker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker(workers, w)
	queues[backchannel.QueueName] = river.QueueConfig{
		MaxWorkers: int(w.config.Workers),
	}
}

type BackChannelLogoutWorkerConfig struct {
	Workers             uint8
	TransactionDuration time.Duration
	MaxAttempts         uint8
	MaxTtl              time.Duration
}
