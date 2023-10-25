package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/member"
)

const (
	MemberUserIDCol         = "user_id"
	MemberRolesCol          = "roles"
	MemberUserResourceOwner = "user_resource_owner"
	MemberUserOwnerRemoved  = "user_owner_removed"

	MemberCreationDate  = "creation_date"
	MemberChangeDate    = "change_date"
	MemberSequence      = "sequence"
	MemberResourceOwner = "resource_owner"
	MemberInstanceID    = "instance_id"
	MemberOwnerRemoved  = "owner_removed"
)

var (
	memberColumns = []*handler.InitColumn{
		handler.NewColumn(MemberCreationDate, handler.ColumnTypeTimestamp),
		handler.NewColumn(MemberChangeDate, handler.ColumnTypeTimestamp),
		handler.NewColumn(MemberUserIDCol, handler.ColumnTypeText),
		handler.NewColumn(MemberUserResourceOwner, handler.ColumnTypeText),
		handler.NewColumn(MemberUserOwnerRemoved, handler.ColumnTypeBool, handler.Default(false)),
		handler.NewColumn(MemberRolesCol, handler.ColumnTypeTextArray, handler.Nullable()),
		handler.NewColumn(MemberSequence, handler.ColumnTypeInt64),
		handler.NewColumn(MemberResourceOwner, handler.ColumnTypeText),
		handler.NewColumn(MemberInstanceID, handler.ColumnTypeText),
		handler.NewColumn(MemberOwnerRemoved, handler.ColumnTypeBool, handler.Default(false)),
	}
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

func reduceMemberAdded(e member.MemberAddedEvent, userResourceOwner string, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		cols: []handler.Column{
			handler.NewCol(MemberUserIDCol, e.UserID),
			handler.NewCol(MemberUserResourceOwner, userResourceOwner),
			handler.NewCol(MemberUserOwnerRemoved, false),
			handler.NewCol(MemberRolesCol, database.TextArray[string](e.Roles)),
			handler.NewCol(MemberCreationDate, e.CreatedAt()),
			handler.NewCol(MemberChangeDate, e.CreatedAt()),
			handler.NewCol(MemberSequence, e.Sequence()),
			handler.NewCol(MemberResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(MemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(MemberOwnerRemoved, false),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.NewCreateStatement(&e, config.cols), nil
}

func reduceMemberChanged(e member.MemberChangedEvent, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		cols: []handler.Column{
			handler.NewCol(MemberRolesCol, database.TextArray[string](e.Roles)),
			handler.NewCol(MemberChangeDate, e.CreatedAt()),
			handler.NewCol(MemberSequence, e.Sequence()),
		},
		conds: []handler.Condition{
			handler.NewCond(MemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(MemberUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.NewUpdateStatement(&e, config.cols, config.conds), nil
}

func reduceMemberCascadeRemoved(e member.MemberCascadeRemovedEvent, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(MemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(MemberUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}
	return handler.NewDeleteStatement(&e, config.conds), nil
}

func reduceMemberRemoved(e eventstore.Event, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(MemberInstanceID, e.Aggregate().InstanceID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}
	return handler.NewDeleteStatement(e, config.conds), nil
}

func multiReduceMemberOwnerRemoved(e eventstore.Event, opts ...reduceMemberOpt) func(eventstore.Event) handler.Exec {
	config := reduceMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(MemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(MemberResourceOwner, e.Aggregate().ID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.AddDeleteStatement(
		config.conds,
	)
}

func memberUserOwnerRemovedConds(e eventstore.Event, opts ...reduceMemberOpt) []handler.Condition {
	config := reduceMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(MemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(MemberUserResourceOwner, e.Aggregate().ID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}
	return config.conds
}

func memberUserOwnerRemovedCols(e eventstore.Event) []handler.Column {
	return []handler.Column{
		handler.NewCol(MemberChangeDate, e.CreatedAt()),
		handler.NewCol(MemberSequence, e.Sequence()),
		handler.NewCol(MemberUserOwnerRemoved, true),
	}
}

func reduceMemberUserOwnerRemoved(e eventstore.Event, opts ...reduceMemberOpt) (*handler.Statement, error) {
	return handler.NewDeleteStatement(
		e,
		memberUserOwnerRemovedConds(e, opts...),
	), nil
}

func multiReduceMemberUserOwnerRemoved(e eventstore.Event, opts ...reduceMemberOpt) func(eventstore.Event) handler.Exec {
	return handler.AddDeleteStatement(
		memberUserOwnerRemovedConds(e, opts...),
	)
}

func setMemberContext(aggregate *eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(context.Background(), aggregate.InstanceID)
}
