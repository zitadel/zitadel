package query

import "github.com/caos/zitadel/internal/query/projection"

var (
	projectMemberTable = table{
		name:  projection.ProjectMemberProjectionTable,
		alias: "members",
	}
	ProjectMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: projectMemberTable,
	}
	ProjectMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: projectMemberTable,
	}
	ProjectMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: projectMemberTable,
	}
	ProjectMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: projectMemberTable,
	}
	ProjectMemberSequence = Column{
		name:  projection.MemberSequence,
		table: projectMemberTable,
	}
	ProjectMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: projectMemberTable,
	}
	ProjectMemberProjectID = Column{
		name:  projection.ProjectMemberProjectIDCol,
		table: projectMemberTable,
	}
)
