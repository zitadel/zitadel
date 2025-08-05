package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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
)

type authNKeyProjection struct{}

func newAuthNKeyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(authNKeyProjection))
}

func (*authNKeyProjection) Name() string {
	return AuthNKeyTable
}

func (*authNKeyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(AuthNKeyIDCol, handler.ColumnTypeText),
			handler.NewColumn(AuthNKeyCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(AuthNKeyChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(AuthNKeyResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(AuthNKeyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(AuthNKeyAggregateIDCol, handler.ColumnTypeText),
			handler.NewColumn(AuthNKeySequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(AuthNKeyObjectIDCol, handler.ColumnTypeText),
			handler.NewColumn(AuthNKeyExpirationCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(AuthNKeyIdentifierCol, handler.ColumnTypeText),
			handler.NewColumn(AuthNKeyPublicKeyCol, handler.ColumnTypeBytes),
			handler.NewColumn(AuthNKeyEnabledCol, handler.ColumnTypeBool, handler.Default(true)),
			handler.NewColumn(AuthNKeyTypeCol, handler.ColumnTypeEnum, handler.Default(0)),
		},
			handler.NewPrimaryKey(AuthNKeyInstanceIDCol, AuthNKeyIDCol),
			handler.WithIndex(handler.NewIndex("enabled", []string{AuthNKeyEnabledCol})),
			handler.WithIndex(handler.NewIndex("identifier", []string{AuthNKeyIdentifierCol})),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{AuthNKeyResourceOwnerCol})),
			handler.WithIndex(handler.NewIndex("creation_date", []string{AuthNKeyCreationDateCol})),
			handler.WithIndex(handler.NewIndex("expiration_date", []string{AuthNKeyExpirationCol})),
		),
	)
}

func (p *authNKeyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Dgb32", "reduce.wrong.event.type %v", []eventstore.EventType{project.ApplicationKeyAddedEventType, user.MachineKeyAddedEventType})
	}
	return handler.NewCreateStatement(
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
			return handler.NewNoOpStatement(event), nil
		}
		appID = e.AppID
		enabled = *e.AuthMethodType == domain.APIAuthMethodTypePrivateKeyJWT
		changeDate = e.CreationDate()
		sequence = e.Sequence()
	case *project.OIDCConfigChangedEvent:
		if e.AuthMethodType == nil {
			return handler.NewNoOpStatement(event), nil
		}
		appID = e.AppID
		enabled = *e.AuthMethodType == domain.OIDCAuthMethodTypePrivateKeyJWT
		changeDate = e.CreationDate()
		sequence = e.Sequence()
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Dbrt1", "reduce.wrong.event.type %v", []eventstore.EventType{project.APIConfigChangedType, project.OIDCConfigChangedType})
	}
	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-BGge42", "reduce.wrong.event.type %v", []eventstore.EventType{project.ApplicationKeyRemovedEventType, project.ApplicationRemovedType, project.ProjectRemovedType, user.MachineKeyRemovedEventType, user.UserRemovedType})
	}
	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Hyd1f", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AuthNKeyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(AuthNKeyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
