package query

import (
	"context"
	"database/sql"
	_ "embed"
	"sync"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

// introspectionTriggerHandlers slice can only be created after zitadel
// is fully initialized, otherwise the handlers are nil.
// OnceValue takes care of creating the slice on the first request
// and than will always return the same slice on subsequent requests.
var introspectionTriggerHandlers = sync.OnceValue(func() []*handler.Handler {
	return append(oidcUserInfoTriggerHandlers(),
		projection.AppProjection,
		projection.OIDCSettingsProjection,
		projection.AuthNKeyProjection,
	)
})

func TriggerIntrospectionProjections(ctx context.Context) {
	triggerBatch(ctx, introspectionTriggerHandlers()...)
}

type IntrospectionClient struct {
	ClientID     string
	ClientSecret *crypto.CryptoValue
	ProjectID    string
	PublicKeys   database.Map[[]byte]
}

//go:embed embed/introspection_client_by_id.sql
var introspectionClientByIDQuery string

func (q *Queries) GetIntrospectionClientByID(ctx context.Context, clientID string, getKeys bool) (_ *IntrospectionClient, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var (
		instanceID = authz.GetInstance(ctx).InstanceID()
		client     = new(IntrospectionClient)
	)

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&client.ClientID, &client.ClientSecret, &client.ProjectID, &client.PublicKeys)
	},
		introspectionClientByIDQuery,
		instanceID, clientID, getKeys,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
