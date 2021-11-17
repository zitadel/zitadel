package projection

import (
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/member"
)

const (
	MemberUserIDCol = "user_id"
	MemberRolesCol  = "roles"

	MemberCreationDate  = "creation_date"
	MemberChangeDate    = "change_date"
	MemberSequence      = "sequence"
	MemberResourceOwner = "resource_owner"
)

func reduceMemberAdded(e member.MemberAddedEvent, aggregateIDCol string) (*handler.Statement, error) {
	return crdb.NewCreateStatement(
		&e,
		[]handler.Column{
			handler.NewCol(aggregateIDCol, e.Aggregate().ResourceOwner),
			handler.NewCol(MemberUserIDCol, e.UserID),
			handler.NewCol(MemberRolesCol, e.Roles),
			handler.NewCol(MemberCreationDate, e.CreationDate()),
			handler.NewCol(MemberChangeDate, e.CreationDate()),
			handler.NewCol(MemberSequence, e.Sequence()),
			handler.NewCol(MemberResourceOwner, e.Aggregate().ResourceOwner),
		},
	), nil
}

func reduceMemberChanged(e member.MemberChangedEvent, aggregateIDCol string) (*handler.Statement, error) {
	return crdb.NewUpdateStatement(
		&e,
		[]handler.Column{
			handler.NewCol(MemberRolesCol, e.Roles),
			handler.NewCol(MemberChangeDate, e.CreationDate()),
			handler.NewCol(MemberSequence, e.Sequence()),
			handler.NewCol(MemberResourceOwner, e.Aggregate().ResourceOwner),
		},
		[]handler.Condition{
			handler.NewCond(aggregateIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(MemberUserIDCol, e.UserID),
		},
	), nil
}

func reduceMemberCascadeRemoved(e member.MemberCascadeRemovedEvent, aggregateIDCol string) (*handler.Statement, error) {
	return crdb.NewDeleteStatement(
		&e,
		[]handler.Condition{
			handler.NewCond(aggregateIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(MemberUserIDCol, e.UserID),
		},
	), nil
}

func reduceMemberRemoved(e member.MemberRemovedEvent, aggregateIDCol string) (*handler.Statement, error) {
	return crdb.NewDeleteStatement(
		&e,
		[]handler.Condition{
			handler.NewCond(aggregateIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(MemberUserIDCol, e.UserID),
		},
	), nil
}
