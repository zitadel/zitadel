package member

import "github.com/caos/zitadel/internal/eventstore/v2"

type Aggregate struct {
	eventstore.Aggregate

	UserID string
	Roles  []string
}

func NewAggregate(aggregate *eventstore.Aggregate, userID string, roles ...string) *Aggregate {
	return &Aggregate{
		Aggregate: *aggregate,
		Roles:     roles,
		UserID:    userID,
	}
}
