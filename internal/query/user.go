package query

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	userTable = table{
		name: projection.UserTable,
	}
	UserIDCol = Column{
		name:  projection.UserIDCol,
		table: userTable,
	}
	UserCreationDateCol = Column{
		name:  projection.UserCreationDateCol,
		table: userTable,
	}
	UserChangeDateCol = Column{
		name:  projection.UserChangeDateCol,
		table: userTable,
	}
	UserResourceOwnerCol = Column{
		name:  projection.UserResourceOwnerCol,
		table: userTable,
	}
	UserStateCol = Column{
		name:  projection.UserStateCol,
		table: userTable,
	}
	UserSequenceCol = Column{
		name:  projection.UserSequenceCol,
		table: userTable,
	}
	UserUsernameCol = Column{
		name:  projection.UserUsernameCol,
		table: userTable,
	}
)

var (
	humanTable = table{
		name: projection.UserHumanTable,
	}
	// profile
	HumanUserIDCol = Column{
		name:  projection.HumanUserIDCol,
		table: humanTable,
	}
	HumanFirstNameCol = Column{
		name:  projection.HumanFirstNameCol,
		table: humanTable,
	}
	HumanLastNameCol = Column{
		name:  projection.HumanLastNameCol,
		table: humanTable,
	}
	HumanNickNameCol = Column{
		name:  projection.HumanNickNameCol,
		table: humanTable,
	}
	HumanDisplayNameCol = Column{
		name:  projection.HumanDisplayNameCol,
		table: humanTable,
	}
	HumanPreferredLanguageCol = Column{
		name:  projection.HumanPreferredLanguageCol,
		table: humanTable,
	}
	HumanGenderCol = Column{
		name:  projection.HumanGenderCol,
		table: humanTable,
	}
	HumanAvaterURLCol = Column{
		name:  projection.HumanAvaterURLCol,
		table: humanTable,
	}

	// email
	HumanEmailCol = Column{
		name:  projection.HumanEmailCol,
		table: humanTable,
	}
	HumanIsEmailVerifiedCol = Column{
		name:  projection.HumanIsEmailVerifiedCol,
		table: humanTable,
	}

	// phone
	HumanPhoneCol = Column{
		name:  projection.HumanPhoneCol,
		table: humanTable,
	}
	HumanIsPhoneVerifiedCol = Column{
		name:  projection.HumanIsPhoneVerifiedCol,
		table: humanTable,
	}
)

var (
	machineTable = table{
		name: projection.UserMachineTable,
	}
	MachineUserIDCol = Column{
		name:  projection.MachineUserIDCol,
		table: machineTable,
	}
	MachineNameCol = Column{
		name:  projection.MachineNameCol,
		table: machineTable,
	}
	MachineDescriptionCol = Column{
		name:  projection.MachineDescriptionCol,
		table: machineTable,
	}
)

func (q *Queries) UserEvents(ctx context.Context, orgID, userID string, sequence uint64) ([]eventstore.EventReader, error) {
	query := NewUserEventSearchQuery(userID, orgID, sequence)
	return q.eventstore.FilterEvents(ctx, query)
}
