package handlers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_enrichCtx(t *testing.T) {
	tests := []struct {
		name       string
		origin     string
		wantOrigin string
		wantErr    bool
	}{
		// e.g. a back-channel logout not triggered by a request
		{name: "empty origin", origin: "", wantOrigin: "://"},
		{name: "https origin", origin: "https://trusted.example.com", wantOrigin: "https://trusted.example.com"},
		{name: "http origin with port", origin: "http://localhost:8080", wantOrigin: "http://localhost:8080"},
		{name: "ipv6 origin", origin: "https://[::1]:8443", wantOrigin: "https://[::1]:8443"},
		{name: "empty port is removed", origin: "https://trusted.example.com:", wantOrigin: "https://trusted.example.com"},
		// stored before the host headers were validated: url.Parse would use "evil.example" as host
		{name: "userinfo with trusted host and port", origin: "https://trusted.example.com:443@evil.example", wantErr: true},
		{name: "userinfo", origin: "https://a:b@evil.com", wantErr: true},
		{name: "invalid host", origin: "https://]", wantErr: true},
		{name: "unsupported scheme", origin: "javascript://trusted.example.com", wantErr: true},
		{name: "opaque", origin: "https:trusted.example.com", wantErr: true},
		{name: "no host", origin: "https://", wantErr: true},
		{name: "unparsable", origin: "https://%zz", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := enrichCtx(context.Background(), tt.origin)
			if tt.wantErr {
				assert.True(t, zerrors.IsErrorInvalidArgument(err), "expected invalid argument error, got: %v", err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantOrigin, http_utils.DomainContext(ctx).Origin())
		})
	}
}

// Test_Origin_invalidStoredOriginFallsBackToPrimaryDomain guards the regression path of the fix:
// an invalid origin stored on an older event (e.g. the userinfo trick that url.Parse would read as
// "evil.example") must not poison the links built from it. Origin falls back to the primary instance
// domain instead of using the stored origin. Without this, reverting the fallback or dropping the
// validation in parseOrigin would bring the poisoned links back with no failing test.
func Test_Origin_invalidStoredOriginFallsBackToPrimaryDomain(t *testing.T) {
	ctrl := gomock.NewController(t)
	queries := mock.NewMockQueries(ctrl)
	queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
		Domains: []*query.InstanceDomain{{
			Domain:    instancePrimaryDomain,
			IsPrimary: true,
		}},
	}, nil)

	n := &NotificationQueries{
		Queries:        queries,
		externalPort:   externalPort,
		externalSecure: externalSecure,
	}

	event := &user.HumanInitialCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
			InstanceID:   "instance",
			AggregateID:  "user",
			CreationDate: time.Now().UTC(),
			Typ:          user.HumanInitialCodeAddedType,
		}),
		TriggeredAtOrigin: "https://trusted.example.com:443@evil.example",
	}

	ctx, err := n.Origin(context.Background(), event)
	require.NoError(t, err)
	assert.Equal(t,
		fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
		http_utils.DomainContext(ctx).Origin(),
	)
}
