package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SamlRequest struct {
	ID           string
	CreationDate time.Time
	LoginClient  string
	Issuer       string
	ACS          string
	RelayState   string
	Binding      string
}

func (a *SamlRequest) checkLoginClient(ctx context.Context, permissionCheck domain.PermissionCheck) error {
	if uid := authz.GetCtxData(ctx).UserID; uid != a.LoginClient {
		return permissionCheck(ctx, domain.PermissionSessionRead, authz.GetInstance(ctx).InstanceID(), "")
	}
	return nil
}

//go:embed saml_request_by_id.sql
var samlRequestByIDQuery string

func (q *Queries) samlRequestByIDQuery(ctx context.Context) string {
	return fmt.Sprintf(samlRequestByIDQuery, q.client.Timetravel(call.Took(ctx)))
}

func (q *Queries) SamlRequestByID(ctx context.Context, shouldTriggerBulk bool, id string, checkLoginClient bool) (_ *SamlRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerSamlRequestProjection")
		ctx, err = projection.SamlRequestProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	dst := new(SamlRequest)
	err = q.client.QueryRowContext(
		ctx,
		func(row *sql.Row) error {
			return row.Scan(
				&dst.ID, &dst.CreationDate, &dst.LoginClient, &dst.Issuer, &dst.ACS, &dst.RelayState, &dst.Binding,
			)
		},
		q.samlRequestByIDQuery(ctx),
		id, authz.GetInstance(ctx).InstanceID(),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-Thee9", "Errors.SamlRequest.NotExisting")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Ou8ue", "Errors.Internal")
	}

	if checkLoginClient {
		if err = dst.checkLoginClient(ctx, q.checkPermission); err != nil {
			return nil, err
		}
	}

	return dst, nil
}
