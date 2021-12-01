package query

import "github.com/caos/zitadel/internal/query/projection"

var (
	iamMemberTable = table{
		name: projection.IAMMemberProjectionTable,
	}
	IAMMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: iamMemberTable,
	}
	IAMMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: iamMemberTable,
	}
	IAMMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: iamMemberTable,
	}
	IAMMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: iamMemberTable,
	}
	IAMMemberSequence = Column{
		name:  projection.MemberSequence,
		table: iamMemberTable,
	}
	IAMMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: iamMemberTable,
	}
	IAMMemberIAMID = Column{
		name:  projection.IAMMemberIAMIDCol,
		table: iamMemberTable,
	}
)
