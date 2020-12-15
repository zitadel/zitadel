package query

import (
	"context"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *QuerySide) IAMMemberByID(ctx context.Context, iamID, userID string) (member *iam_repo.MemberReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	member = iam_repo.NewMemberReadModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}
