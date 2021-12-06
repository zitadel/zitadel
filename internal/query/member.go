package query

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type Members struct {
	SearchResponse
	Members []*Member
}

type Member struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	UserID             string
	Roles              []string
	PreferredLoginName string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	AvatarURL          string
}

func (r *Queries) IAMMemberByID(ctx context.Context, iamID, userID string) (member *IAMMemberReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	member = NewIAMMemberReadModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}
