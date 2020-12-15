package query

import (
	"context"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *QuerySide) IAMMemberByID(ctx context.Context, iamID, userID string) (member *IAMMemberReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	member = NewIAMMemberReadModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}
