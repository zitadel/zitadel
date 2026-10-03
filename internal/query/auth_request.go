package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"strings"
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
	// PrivateLabelingOrgID is the organization whose branding should be shown for this
	// request, resolved from the org scope, the requested project's private labeling
	// setting, or the instance default organization.
	PrivateLabelingOrgID string
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
	dst.PrivateLabelingOrgID = q.resolvePrivateLabelingOrgID(ctx, dst.ClientID, dst.Scope)

	if checkLoginClient {
		if err = dst.checkLoginClient(ctx, q.checkPermission); err != nil {
			return nil, err
		}
	}

	return dst, nil
}

// resolvePrivateLabelingOrgID determines which organization's branding (private labeling)
// should be shown for an auth request, reusing the same rule as the legacy login
// (domain.AuthRequest.PrivateLabelingOrgID): an explicit organization scope wins, otherwise
// the requested project's private labeling setting decides, otherwise the instance default
// organization. The user is not known at this point, so the user-org branch never applies here.
// Branding must never fail the request, so any lookup error falls back to the default org.
func (q *Queries) resolvePrivateLabelingOrgID(ctx context.Context, clientID string, scopes []string) string {
	ar := new(domain.AuthRequest)
	for _, scope := range scopes {
		if orgID, ok := strings.CutPrefix(scope, domain.OrgIDScope); ok {
			ar.RequestedOrgID = orgID
			break
		}
	}
	if ar.RequestedOrgID == "" && clientID != "" {
		if project, err := q.ProjectByClientID(ctx, clientID); err == nil && project != nil {
			ar.PrivateLabelingSetting = project.PrivateLabelingSetting
			ar.ApplicationResourceOwner = project.ResourceOwner
		}
	}
	return ar.PrivateLabelingOrgID(authz.GetInstance(ctx).DefaultOrganisationID())
}
