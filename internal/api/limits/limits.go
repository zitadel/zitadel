package limits

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type limitsKey struct{}

var key = (*limitsKey)(nil)

type Loader struct {
	querier Querier
}

// NewLoader makes it easy to deduplicate multiple limit queries while handling a single request.
// The returned limitsLoader itself is stateless.
// Once the limits are loaded, they are attached to the context.
// Within a contexts lifetime, the Loader also guarantees that even in error cases, the limits are only queried once.
// Therefore, there won't be any circular calls as long as the passed context is a child of a previously passed context.
func NewLoader(querier Querier) *Loader {
	return &Loader{querier}
}

// Querier abstracts query.Queries.Limits to avoid circular dependencies.
type Querier interface {
	Limits(ctx context.Context, resourceOwner string) (limits *Limits, err error)
}

type Limits struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	AuditLogRetention *time.Duration
	Block             *bool
}

// Load ensures that if limits are already attached to the context, they are not queried again.
// Use the returned context for further calls to Load.
func (l *Loader) Load(ctx context.Context, instanceID string) (context.Context, Limits) {
	ctxLimits, ok := ctx.Value(key).(*Limits)
	if ok {
		return ctx, *ctxLimits
	}
	queriedLimits, err := l.querier.Limits(ctx, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		logging.WithFields("instance id", instanceID).OnError(err).Error("unable to load limits")
	}
	if queriedLimits == nil {
		queriedLimits = &Limits{}
	}
	return context.WithValue(ctx, key, queriedLimits), *queriedLimits
}
