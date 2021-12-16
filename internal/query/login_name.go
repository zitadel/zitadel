package query

import "github.com/caos/zitadel/internal/query/projection"

var (
	loginNameTable = table{
		name: projection.LoginNameProjectionTable,
	}
	LoginNameUserIDCol = Column{
		name:  "user_id",
		table: loginNameTable,
	}
	LoginNameNameCol = Column{
		name:  projection.LoginNameCol,
		table: loginNameTable,
	}
	LoginNameIsPrimaryCol = Column{
		name:  projection.LoginNameDomainIsPrimaryCol,
		table: loginNameTable,
	}
)
