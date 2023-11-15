package query

import "github.com/zitadel/zitadel/v2/internal/query/projection"

var (
	loginNameTable = table{
		name:          projection.LoginNameProjectionTable,
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
	LoginNameOwnerRemovedDomainCol = Column{
		name:  projection.LoginNameOwnerRemovedDomainCol,
		table: loginNameTable,
	}
	LoginNameOwnerRemovedUserCol = Column{
		name:  projection.LoginNameOwnerRemovedUserCol,
		table: loginNameTable,
	}
	LoginNameOwnerRemovedPolicyCol = Column{
		name:  projection.LoginNameOwnerRemovedPolicyCol,
		table: loginNameTable,
	}
)

func addLoginNameWithoutOwnerRemoved(eq map[string]interface{}) {
	eq[LoginNameOwnerRemovedDomainCol.identifier()] = false
	eq[LoginNameOwnerRemovedUserCol.identifier()] = false
	eq[LoginNameOwnerRemovedPolicyCol.identifier()] = false
}
