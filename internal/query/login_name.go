package query

import "github.com/zitadel/zitadel/internal/query/projection"

var (
	loginNameTable = table{
		name:          "projections.login_names4",
		instanceIDCol: projection.LoginNameUserInstanceIDCol,
	}
	LoginNameUserIDCol = Column{
		name:  projection.LoginNameUserCol,
		table: loginNameTable,
	}
	LoginNameNameCol = Column{
		name:  projection.LoginNameCol,
		table: loginNameTable,
	}
	LoginNameIsPrimaryCol = Column{
		name:  projection.LoginNameIsPrimaryCol,
		table: loginNameTable,
	}
	LoginNameInstanceIDCol = Column{
		name:  projection.LoginNameInstanceIDCol,
		table: loginNameTable,
	}
	LoginNameResourceOwnerCol = Column{
		name:  "resource_owner",
		table: loginNameTable,
	}
)
