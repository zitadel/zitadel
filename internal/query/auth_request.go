package query

import (
	"context"
	"database/sql"
	_ "embed"
	errs "errors"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type AuthRequest struct {
	ID           string
	CreationDate time.Time
	ClientID     string
	Scope        []string
	RedirectURI  string
	Prompt       []domain.Prompt
	UiLocales    []string
	LoginHint    *string
	MaxAge       *time.Duration
	HintUserID   *string
}

//go:embed embed/auth_request_by_id.sql
var authRequestByIDQuery string

func (q *Queries) AuthRequestByID(ctx context.Context, shouldTriggerBulk bool, id string) (_ *AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.AuthRequestProjection.Trigger(ctx)
	}

	var (
		scope   database.StringArray
		prompt  database.EnumArray[domain.Prompt]
		locales database.StringArray
	)

	dst := new(AuthRequest)
	err = q.client.DB.QueryRowContext(ctx, authRequestByIDQuery, id).Scan(
		&dst.ID, &dst.CreationDate, &dst.ClientID, &scope, &dst.RedirectURI,
		&prompt, &locales, &dst.LoginHint, &dst.MaxAge, &dst.HintUserID,
	)
	if errs.Is(err, sql.ErrNoRows) {
		return nil, errors.ThrowNotFound(err, "QUERY-Thee9", "Errors.AuthRequest.NotExisting")
	}
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Ou8ue", "Errors.Internal")
	}

	dst.Scope = scope
	dst.Prompt = prompt
	dst.UiLocales = locales

	return dst, nil
}
