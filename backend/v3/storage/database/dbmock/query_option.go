package dbmock

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// QueryOptions converts database.QueryOption to *database.QueryOpts which
// can then be used as a gomock Matcher.
func QueryOptions(optsFunc database.QueryOption) *database.QueryOpts {
	opts := &database.QueryOpts{}

	optsFunc(opts)

	return opts
}
