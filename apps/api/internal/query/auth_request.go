package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AuthRequest struct {
	ID           string
	CreationDate time.Time
	LoginClient  string
	ClientID     string
	Scope        []string
	RedirectURI  string
	Prompt       []domain.Prompt
	UiLocales    []string
	LoginHint    *string
	MaxAge       *time.Duration
	HintUserID   *string
}

func (a *AuthRequest) checkLoginClient(ctx context.Context, permissionCheck domain.PermissionCheck) error {
	if uid := authz.GetCtxData(ctx).UserID; uid != a.LoginClient {
		return permissionCheck(ctx, domain.PermissionSessionRead, authz.GetInstance(ctx).InstanceID(), "")
	}
	return nil
}

//go:embed auth_request_by_id.sql
var authRequestByIDQuery string

func (q *Queries) AuthRequestByID(ctx context.Context, shouldTriggerBulk bool, id string, checkLoginClient bool) (_ *AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerAuthRequestProjection")
		ctx, err = projection.AuthRequestProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	var (
		scope   database.TextArray[string]
		prompt  database.NumberArray[domain.Prompt]
		locales database.TextArray[string]
	)

	dst := new(AuthRequest)
	err = q.client.QueryRowContext(
		ctx,
		func(row *sql.Row) error {
			return row.Scan(
				&dst.ID, &dst.CreationDate, &dst.LoginClient, &dst.ClientID, &scope, &dst.RedirectURI,
				&prompt, &locales, &dst.LoginHint, &dst.MaxAge, &dst.HintUserID,
			)
		},
		authRequestByIDQuery,
		id, authz.GetInstance(ctx).InstanceID(),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-Thee9", "Errors.AuthRequest.NotExisting")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Ou8ue", "Errors.Internal")
	}

	dst.Scope = scope
	dst.Prompt = prompt
	dst.UiLocales = locales

	if checkLoginClient {
		if err = dst.checkLoginClient(ctx, q.checkPermission); err != nil {
			return nil, err
		}
	}

	return dst, nil
}
