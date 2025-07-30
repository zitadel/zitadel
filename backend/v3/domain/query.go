package domain

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type Pagination struct {
	Limit        uint32
	Offset       uint32
	Ascending    bool
	OrderColumns database.Columns
}

// applyOnOrgsQuery implements OrgsQueryOpts.
func (p Pagination) applyOnOrgsQuery(query *OrgsQuery) {
	query.pagination = p
}

var _ OrgsQueryOpts = (*Pagination)(nil)
