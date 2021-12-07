package query

import "github.com/caos/zitadel/internal/query/projection"

var (
	orgMemberTable = table{
		name:  projection.OrgMemberProjectionTable,
		alias: "m",
	}
	OrgMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: orgMemberTable,
	}
	OrgMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: orgMemberTable,
	}
	OrgMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: orgMemberTable,
	}
	OrgMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: orgMemberTable,
	}
	OrgMemberSequence = Column{
		name:  projection.MemberSequence,
		table: orgMemberTable,
	}
	OrgMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: orgMemberTable,
	}
	OrgMemberOrgID = Column{
		name:  projection.OrgMemberOrgIDCol,
		table: orgMemberTable,
	}
)
