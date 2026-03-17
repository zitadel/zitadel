package instrumentation

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/rs/xid"

	"github.com/zitadel/zitadel/internal/api/call"
)

// allow injection for testing
var xidWithTime = xid.NewWithTime

// WithRequestDetails creates a new context with request details, including a unique request ID.
// The request ID is generated using the call timestamp as returned by [call.FromContext].
//
// Note: The returned context carries a mutable pointer to requestDetails.
// This allows subsequent calls to SetInstanceID/SetUserID to update the details
// visible to the logger, even if the logger holds the original context created here.
func WithRequestDetails(ctx context.Context, instanceHost, publicHost string) context.Context {
	details := &requestDetails{
		id:           xidWithTime(call.FromContext(ctx)),
		instanceHost: instanceHost,
		publicHost:   publicHost,
	}
	return context.WithValue(ctx, ctxKey{}, details)
}

func getRequestDetails(ctx context.Context) (*requestDetails, bool) {
	details, ok := ctx.Value(ctxKey{}).(*requestDetails)
	return details, ok
}

// SetInstanceID sets the instance ID in the request details stored in the context.
// It uses a mutex to update the shared state safely, ensuring upstream loggers see the change.
func SetInstanceID(ctx context.Context, instanceID string) {
	if details, ok := getRequestDetails(ctx); ok {
		details.mtx.Lock()
		details.instanceID = instanceID
		details.mtx.Unlock()
	}
}

// SetUserID sets the user ID in the request details stored in the context.
// It uses a mutex to update the shared state safely, ensuring upstream loggers see the change.
func SetUserID(ctx context.Context, userID string) {
	if details, ok := getRequestDetails(ctx); ok {
		details.mtx.Lock()
		details.userID = userID
		details.mtx.Unlock()
	}
}

// GetRequestID retrieves the request ID from the context.
// [xid.NilID] is returned if no request details are found in the context.
func GetRequestID(ctx context.Context) xid.ID {
	d, ok := getRequestDetails(ctx)
	if !ok {
		return xid.NilID()
	}
	d.mtx.Lock()
	defer d.mtx.Unlock()
	return d.id
}

type ctxKey struct{}

type requestDetails struct {
	mtx          sync.Mutex
	id           xid.ID
	instanceHost string
	publicHost   string // may optionally be set through header (login v2)
	instanceID   string
	userID       string
}

func (d *requestDetails) slogAttributes() []any {
	attributes := make([]any, 0, 5)

	d.mtx.Lock()
	defer d.mtx.Unlock()

	attributes = append(attributes,
		slog.String("id", d.id.String()),
		slog.String("instance_host", d.instanceHost),
	)
	if d.publicHost != "" {
		attributes = append(attributes, slog.String("public_host", d.publicHost))
	}
	if d.instanceID != "" {
		attributes = append(attributes, slog.String("instance_id", d.instanceID))
	}
	if d.userID != "" {
		attributes = append(attributes, slog.String("user_id", d.userID))
	}

	return attributes
}

func requestDetailsExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	if d, ok := getRequestDetails(ctx); ok {
		return []slog.Attr{
			slog.Group("request", d.slogAttributes()...),
		}
	}
	return nil
}
