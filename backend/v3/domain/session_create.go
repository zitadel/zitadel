package domain

import (
	"context"
	"net"
	"net/http"
	"strings"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type CreateSessionCommand struct {
	InstanceID string `json:"instance_id"`
	UserAgent  *session_grpc.UserAgent
	Metas      map[string][]byte
	Lifetime   *durationpb.Duration

	idGenerator id.Generator

	SessionID *string
}

func NewCreateSessionCommand(instanceID string, userAgent *session_grpc.UserAgent, metas map[string][]byte, lifetime *durationpb.Duration, idGenerator id.Generator) *CreateSessionCommand {
	idGen := id.SonyFlakeGenerator()
	if idGenerator != nil {
		idGen = idGenerator
	}

	return &CreateSessionCommand{
		InstanceID:  strings.TrimSpace(instanceID),
		UserAgent:   userAgent,
		idGenerator: idGen,
		Metas:       metas,
		Lifetime:    lifetime,
	}
}

// RequiresTransaction implements [Transactional].
func (c *CreateSessionCommand) RequiresTransaction() {}

// Events implements [Commander].
func (c *CreateSessionCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	var oldUA *old_domain.UserAgent
	ua := c.getUserAgent()

	if ua != nil {
		oldUA = &old_domain.UserAgent{
			FingerprintID: ua.FingerprintID,
			IP:            ua.IP,
			Description:   ua.Description,
			Header:        ua.Header,
		}
	}

	sessionAgg := &session.NewAggregate(*c.SessionID, c.InstanceID).Aggregate

	toReturn := make([]eventstore.Command, 1, 3)
	toReturn[0] = session.NewAddedEvent(ctx, sessionAgg, oldUA)

	if c.Metas != nil {
		toReturn = append(toReturn, session.NewMetadataSetEvent(ctx, sessionAgg, c.Metas))
	}

	if c.Lifetime != nil {
		toReturn = append(toReturn, session.NewLifetimeSetEvent(ctx, sessionAgg, c.Lifetime.AsDuration()))
	}

	return toReturn, nil
}

// Execute implements [Commander].
func (c *CreateSessionCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	sessionRepo := opts.sessionRepo

	sessionID, err := c.idGenerator.Next()
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-ngXOIK", "failed generating session ID")
	}

	session := &Session{
		ID:         sessionID,
		InstanceID: c.InstanceID,
		UserAgent:  c.getUserAgent(),
		Metadata:   c.getMetas(sessionID),
		Lifetime:   c.Lifetime.AsDuration(),
	}

	err = sessionRepo.Create(ctx, opts.DB(), session)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-HYKAgF", "failed creating session")
	}

	c.SessionID = &sessionID

	return nil
}

func (c *CreateSessionCommand) getMetas(sessionID string) []SessionMetadata {
	if c.Metas == nil {
		return nil
	}
	toReturn := make([]SessionMetadata, len(c.Metas))
	i := 0

	for metaK, metaV := range c.Metas {
		sm := SessionMetadata{
			Metadata: Metadata{
				InstanceID: c.InstanceID,
				Key:        metaK,
				Value:      metaV,
			},
			SessionID: sessionID,
		}
		toReturn[i] = sm
		i += 1
	}

	return toReturn
}

func (c *CreateSessionCommand) getUserAgent() *SessionUserAgent {
	if c.UserAgent == nil {
		return nil
	}

	toReturn := &SessionUserAgent{
		FingerprintID: c.UserAgent.FingerprintId,
		InstanceID:    c.InstanceID,
		Description:   c.UserAgent.Description,
		IP:            net.ParseIP(c.UserAgent.GetIp()),
	}

	if headers := c.UserAgent.GetHeader(); len(headers) > 0 {
		toReturn.Header = make(http.Header, len(headers))
		for hKey, hValue := range headers {
			toReturn.Header[hKey] = hValue.GetValues()
		}
	}

	return toReturn
}

// String implements [Commander].
func (c *CreateSessionCommand) String() string {
	return "CreateSessionCommand"
}

// Validate implements [Commander].
func (c *CreateSessionCommand) Validate(ctx context.Context, _ *InvokeOpts) (err error) {
	if c.InstanceID == "" {
		c.InstanceID = authz.GetInstance(ctx).InstanceID()
	}

	if c.Lifetime == nil {
		return nil
	}
	asDuration := c.Lifetime.AsDuration()
	if asDuration < 0 {
		return zerrors.ThrowInvalidArgument(nil, "DOM-XA5OMq", "Errors.Session.PositiveLifetime")
	}

	return nil
}

var _ Transactional = (*CreateSessionCommand)(nil)
var _ Commander = (*CreateSessionCommand)(nil)
