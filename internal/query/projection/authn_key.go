package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	AuthNKeyTable            = "projections.authn_keys2"
	AuthNKeyIDCol            = "id"
	AuthNKeyCreationDateCol  = "creation_date"
	AuthNKeyChangeDateCol    = "change_date"
	AuthNKeyResourceOwnerCol = "resource_owner"
	AuthNKeyInstanceIDCol    = "instance_id"
	AuthNKeyAggregateIDCol   = "aggregate_id"
	AuthNKeySequenceCol      = "sequence"
	AuthNKeyObjectIDCol      = "object_id"
	AuthNKeyExpirationCol    = "expiration"
	AuthNKeyIdentifierCol    = "identifier"
	AuthNKeyPublicKeyCol     = "public_key"
	AuthNKeyTypeCol          = "type"
	AuthNKeyEnabledCol       = "enabled"
	AuthNKeyOwnerRemovedCol  = "owner_removed"
)

type authNKeyProjection struct {
	crdb.StatementHandler
}

func newAuthNKeyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *authNKeyProjection {
	p := new(authNKeyProjection)
	config.ProjectionName = AuthNKeyTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(AuthNKeyIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AuthNKeyCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AuthNKeyChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AuthNKeyResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(AuthNKeyInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AuthNKeyAggregateIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AuthNKeySequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(AuthNKeyObjectIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AuthNKeyExpirationCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AuthNKeyIdentifierCol, crdb.ColumnTypeText),
			crdb.NewColumn(AuthNKeyPublicKeyCol, crdb.ColumnTypeBytes),
			crdb.NewColumn(AuthNKeyEnabledCol, crdb.ColumnTypeBool, crdb.Default(true)),
			crdb.NewColumn(AuthNKeyTypeCol, crdb.ColumnTypeEnum, crdb.Default(0)),
			crdb.NewColumn(AuthNKeyOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(AuthNKeyInstanceIDCol, AuthNKeyIDCol),
			crdb.WithIndex(crdb.NewIndex("enabled", []string{AuthNKeyEnabledCol})),
			crdb.WithIndex(crdb.NewIndex("identifier", []string{AuthNKeyIdentifierCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{AuthNKeyOwnerRemovedCol})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *authNKeyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.ApplicationKeyAddedEventType,
					Reduce: p.reduceAuthNKeyAdded,
				},
				{
					Event:  project.ApplicationKeyRemovedEventType,
					Reduce: p.reduceAuthNKeyRemoved,
				},
				{
					Event:  project.APIConfigChangedType,
					Reduce: p.reduceAuthNKeyEnabledChanged,
				},
				{
					Event:  project.OIDCConfigChangedType,
					Reduce: p.reduceAuthNKeyEnabledChanged,
				},
				{
					Event:  project.ApplicationRemovedType,
					Reduce: p.reduceAuthNKeyRemoved,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceAuthNKeyRemoved,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.MachineKeyAddedEventType,
					Reduce: p.reduceAuthNKeyAdded,
				},
				{
					Event:  user.MachineKeyRemovedEventType,
					Reduce: p.reduceAuthNKeyRemoved,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceAuthNKeyRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(AuthNKeyInstanceIDCol),
				},
			},
		},
	}
}

func (p *authNKeyProjection) reduceAuthNKeyAdded(event eventstore.Event) (*handler.Statement, error) {
	var authNKeyEvent struct {
		eventstore.BaseEvent
		keyID      string
		objectID   string
		expiration time.Time
		identifier string
		publicKey  []byte
		keyType    domain.AuthNKeyType
	}
	switch e := event.(type) {
	case *project.ApplicationKeyAddedEvent:
		authNKeyEvent.BaseEvent = e.BaseEvent
		authNKeyEvent.keyID = e.KeyID
		authNKeyEvent.objectID = e.AppID
		authNKeyEvent.expiration = e.ExpirationDate
		authNKeyEvent.identifier = e.ClientID
		authNKeyEvent.publicKey = e.PublicKey
		authNKeyEvent.keyType = e.KeyType
	case *user.MachineKeyAddedEvent:
		authNKeyEvent.BaseEvent = e.BaseEvent
		authNKeyEvent.keyID = e.KeyID
		authNKeyEvent.objectID = e.Aggregate().ID
		authNKeyEvent.expiration = e.ExpirationDate
		authNKeyEvent.identifier = e.Aggregate().ID
		authNKeyEvent.publicKey = e.PublicKey
		authNKeyEvent.keyType = e.KeyType
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Dgb32", "reduce.wrong.event.type %v", []eventstore.EventType{project.ApplicationKeyAddedEventType, user.MachineKeyAddedEventType})
	}
	return crdb.NewCreateStatement(
		&authNKeyEvent,
		[]handler.Column{
			handler.NewCol(AuthNKeyIDCol, authNKeyEvent.keyID),
			handler.NewCol(AuthNKeyCreationDateCol, authNKeyEvent.CreationDate()),
			handler.NewCol(AuthNKeyChangeDateCol, authNKeyEvent.CreationDate()),
			handler.NewCol(AuthNKeyResourceOwnerCol, authNKeyEvent.Aggregate().ResourceOwner),
			handler.NewCol(AuthNKeyInstanceIDCol, authNKeyEvent.Aggregate().InstanceID),
			handler.NewCol(AuthNKeyAggregateIDCol, authNKeyEvent.Aggregate().ID),
			handler.NewCol(AuthNKeySequenceCol, authNKeyEvent.Sequence()),
			handler.NewCol(AuthNKeyObjectIDCol, authNKeyEvent.objectID),
			handler.NewCol(AuthNKeyExpirationCol, authNKeyEvent.expiration),
			handler.NewCol(AuthNKeyIdentifierCol, authNKeyEvent.identifier),
			handler.NewCol(AuthNKeyPublicKeyCol, authNKeyEvent.publicKey),
			handler.NewCol(AuthNKeyTypeCol, authNKeyEvent.keyType),
		},
	), nil
}

func (p *authNKeyProjection) reduceAuthNKeyEnabledChanged(event eventstore.Event) (*handler.Statement, error) {
	var appID string
	var enabled bool
	var changeDate time.Time
	var sequence uint64
	switch e := event.(type) {
	case *project.APIConfigChangedEvent:
		if e.AuthMethodType == nil {
			return crdb.NewNoOpStatement(event), nil
		}
		appID = e.AppID
		enabled = *e.AuthMethodType == domain.APIAuthMethodTypePrivateKeyJWT
		changeDate = e.CreationDate()
		sequence = e.Sequence()
	case *project.OIDCConfigChangedEvent:
		if e.AuthMethodType == nil {
			return crdb.NewNoOpStatement(event), nil
		}
		appID = e.AppID
		enabled = *e.AuthMethodType == domain.OIDCAuthMethodTypePrivateKeyJWT
		changeDate = e.CreationDate()
		sequence = e.Sequence()
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Dbrt1", "reduce.wrong.event.type %v", []eventstore.EventType{project.APIConfigChangedType, project.OIDCConfigChangedType})
	}
	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(AuthNKeyChangeDateCol, changeDate),
			handler.NewCol(AuthNKeySequenceCol, sequence),
			handler.NewCol(AuthNKeyEnabledCol, enabled),
		},
		[]handler.Condition{
			handler.NewCond(AuthNKeyObjectIDCol, appID),
			handler.NewCond(AuthNKeyInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *authNKeyProjection) reduceAuthNKeyRemoved(event eventstore.Event) (*handler.Statement, error) {
	var condition handler.Condition
	switch e := event.(type) {
	case *project.ApplicationKeyRemovedEvent:
		condition = handler.NewCond(AuthNKeyIDCol, e.KeyID)
	case *project.ApplicationRemovedEvent:
		condition = handler.NewCond(AuthNKeyObjectIDCol, e.AppID)
	case *project.ProjectRemovedEvent:
		condition = handler.NewCond(AuthNKeyAggregateIDCol, e.Aggregate().ID)
	case *user.MachineKeyRemovedEvent:
		condition = handler.NewCond(AuthNKeyIDCol, e.KeyID)
	case *user.UserRemovedEvent:
		condition = handler.NewCond(AuthNKeyAggregateIDCol, e.Aggregate().ID)
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-BGge42", "reduce.wrong.event.type %v", []eventstore.EventType{project.ApplicationKeyRemovedEventType, project.ApplicationRemovedType, project.ProjectRemovedType, user.MachineKeyRemovedEventType, user.UserRemovedType})
	}
	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			condition,
			handler.NewCond(AuthNKeyInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *authNKeyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Hyd1f", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AuthNKeyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(AuthNKeyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
