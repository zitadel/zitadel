package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Groups struct {
	SearchResponse
	Groups []*Group
}

type Group struct {
	ID            string
	Name          string
	Description   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
}

// SearchGroups returns the list of groups that match the search criteria
func (q *Queries) SearchGroups(ctx context.Context) (*Groups, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "QUERY-grpfli", "Not implemented")
}
