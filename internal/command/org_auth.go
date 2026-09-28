package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// checkOrgNotDeactivatedAfter returns an invalid refresh token error if the organization
// was deactivated after the given time. This keeps refresh grants dead after reactivation.
func (c *Commands) checkOrgNotDeactivatedAfter(ctx context.Context, orgID string, after time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if orgID == "" || after.IsZero() {
		return nil
	}
	model := &orgDeactivatedAfterModel{
		orgID: orgID,
		after: after,
	}
	if err = c.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return zerrors.ThrowPreconditionFailed(err, "OIDCS-oR9nD", "Errors.Internal")
	}
	if model.deactivated {
		return zerrors.ThrowPreconditionFailed(nil, "OIDCS-oR9nR", "Errors.OIDCSession.RefreshTokenInvalid")
	}
	return nil
}

type orgDeactivatedAfterModel struct {
	orgID string
	after time.Time

	events      int
	deactivated bool
}

func (m *orgDeactivatedAfterModel) Reduce() error {
	m.deactivated = m.events > 0
	return nil
}

func (m *orgDeactivatedAfterModel) AppendEvents(events ...eventstore.Event) {
	m.events += len(events)
}

func (m *orgDeactivatedAfterModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		CreationDateAfter(m.after).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(m.orgID).
		EventTypes(org.OrgDeactivatedEventType).
		Builder()
}
