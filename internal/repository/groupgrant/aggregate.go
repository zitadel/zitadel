// Package groupgrant implements project-role grants for user groups.
//
// Group grants are resolved at query time: tokens, userinfo, and authorization
// queries merge them with personal user grants through the user's memberships
// (see internal/query/userinfo_by_id.sql). Grants are deliberately NOT fanned
// out into per-member user grants — materializing them would turn one
// membership change into N grant events, and group deletion into N
// revocations, recreating the orphaned-permission problem. With query-time
// resolution, revocation is free (the join stops matching) and every derived
// role keeps its provenance (the supplying group).
package groupgrant

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "groupgrant"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id, resourceOwner string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: resourceOwner,
		},
	}
}
