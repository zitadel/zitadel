package member

import "github.com/caos/zitadel/internal/eventstore/v2"

type Aggregate struct {
	eventstore.Aggregate

	UserID string
	Roles  []string
}

func NewMemberAggregate(userID string) *ReadModel {
	return &ReadModel{
		ReadModel: *eventstore.NewReadModel(),
	}
}
