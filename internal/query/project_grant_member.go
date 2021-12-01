package query

import "github.com/caos/zitadel/internal/query/projection"

var (
	projectGrantMemberTable = table{
		name: projection.ProjectGrantMemberProjectionTable,
	}
	ProjectGrantMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberSequence = Column{
		name:  projection.MemberSequence,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberProjectID = Column{
		name:  projection.ProjectGrantMemberProjectIDCol,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberGrantID = Column{
		name:  projection.ProjectGrantMemberGrantIDCol,
		table: projectGrantMemberTable,
	}
)
