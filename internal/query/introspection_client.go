package query

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/jackc/pgtype"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
)

type IntrospectionClient struct {
	ClientID     string
	ClientSecret *crypto.CryptoValue
	ProjectID    string
	PublicKeys   database.Map[[]byte]
}

//go:embed embed/introspection_client_by_id.sql
var introspectionClientByIDQuery string

func (q *Queries) GetIntrospectionClientByID(ctx context.Context, clientID string, getKeys bool) (_ *IntrospectionClient, err error) {
	var (
		instanceID = authz.GetInstance(ctx).InstanceID()
		client     = new(IntrospectionClient)
	)

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		var publicKeys pgtype.ByteaArray
		if err := row.Scan(&client.ClientID, &client.ClientSecret, &client.ProjectID, &publicKeys); err != nil {
			return err
		}
		return publicKeys.AssignTo(&client.PublicKeys)
	},
		introspectionClientByIDQuery,
		instanceID, clientID, getKeys,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
