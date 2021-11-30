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

type reduceMemberConfig struct {
	cols  []handler.Column
	conds []handler.Condition
}

type reduceMemberOpt func(reduceMemberConfig) reduceMemberConfig

func withMemberCol(col string, value interface{}) reduceMemberOpt {
	return func(opt reduceMemberConfig) reduceMemberConfig {
		opt.cols = append(opt.cols, handler.NewCol(col, value))
		return opt
	}
}

func withMemberCond(cond string, value interface{}) reduceMemberOpt {
	return func(opt reduceMemberConfig) reduceMemberConfig {
		opt.conds = append(opt.conds, handler.NewCond(cond, value))
		return opt
	}
}

func reduceMemberAdded(e member.MemberAddedEvent, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		cols: []handler.Column{
			handler.NewCol(MemberUserIDCol, e.UserID),
			handler.NewCol(MemberRolesCol, e.Roles),
			handler.NewCol(MemberCreationDate, e.CreationDate()),
			handler.NewCol(MemberChangeDate, e.CreationDate()),
			handler.NewCol(MemberSequence, e.Sequence()),
			handler.NewCol(MemberResourceOwner, e.Aggregate().ResourceOwner),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return crdb.NewCreateStatement(&e, config.cols), nil
}

func reduceMemberChanged(e member.MemberChangedEvent, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		cols: []handler.Column{
			handler.NewCol(MemberRolesCol, e.Roles),
			handler.NewCol(MemberChangeDate, e.CreationDate()),
			handler.NewCol(MemberSequence, e.Sequence()),
		},
		conds: []handler.Condition{
			handler.NewCond(MemberUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return crdb.NewUpdateStatement(&e, config.cols, config.conds), nil
}

func reduceMemberCascadeRemoved(e member.MemberCascadeRemovedEvent, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(MemberUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}
	return crdb.NewDeleteStatement(&e, config.conds), nil
}

func reduceMemberRemoved(e member.MemberRemovedEvent, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(MemberUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}
	return crdb.NewDeleteStatement(&e, config.conds), nil
}
