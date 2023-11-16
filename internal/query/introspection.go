package query

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var introspectionTriggerHandlers = append(oidcUserInfoTriggerHandlers,
	projection.AppProjection,
	projection.OIDCSettingsProjection,
	projection.AuthNKeyProjection,
)

func TriggerIntrospectionProjections(ctx context.Context) {
	triggerBatch(ctx, introspectionTriggerHandlers...)
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
