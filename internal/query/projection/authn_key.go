package projection

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/internal/repository/user"

	"github.com/caos/zitadel/internal/repository/project"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
)

const (
	AuthNKeyTable            = "zitadel.projections.authn_keys"
	AuthNKeyIDCol            = "id"
	AuthNKeyCreationDateCol  = "creation_date"
	AuthNKeyChangeDateCol    = "change_date"
	AuthNKeyResourceOwnerCol = "resource_owner"
	AuthNKeyAggregateIDCol   = "aggregate_id"
	AuthNKeySequenceCol      = "sequence"
	AuthNKeyObjectIDCol      = "object_id"
	AuthNKeyExpirationCol    = "expiration"
)

const (
	AuthNKeyPublicIDCol         = "key_id"
	AuthNKeyPublicIdentifierCol = "identifier"
	AuthNKeyPublicKeyCol        = "key"
)

type AuthNKeyProjection struct {
	crdb.StatementHandler
}

func NewAuthNKeyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *AuthNKeyProjection {
	p := &AuthNKeyProjection{}
	config.ProjectionName = AuthNKeyTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *AuthNKeyProjection) reducers() []handler.AggregateReducer {
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
					Reduce: p.reduceAuthNKeyRemoved,
				},
				{
					Event:  project.OIDCConfigChangedType,
					Reduce: p.reduceAuthNKeyRemoved,
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
	}
}

func (p *AuthNKeyProjection) reduceAuthNKeyAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var baseEvent eventstore.BaseEvent
	var keyID, objectID, identifier string
	var expiration time.Time
	var publicKey []byte
	switch e := event.(type) {
	case *project.ApplicationKeyAddedEvent:
		baseEvent = e.BaseEvent
		keyID = e.KeyID
		objectID = e.AppID
		expiration = e.ExpirationDate
		identifier = e.ClientID
		publicKey = e.PublicKey
	case *user.MachineKeyAddedEvent:
		baseEvent = e.BaseEvent
		keyID = e.KeyID
		objectID = e.Aggregate().ID
		expiration = e.ExpirationDate
		identifier = e.Aggregate().ID
		publicKey = e.PublicKey
	default:
		logging.LogWithFields("PROJE-Dbr3g", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{project.ApplicationKeyAddedEventType, user.MachineKeyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Dgb32", "reduce.wrong.event.type")
	}
	return crdb.NewMultiStatement(
		&baseEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AuthNKeyIDCol, keyID),
				handler.NewCol(AuthNKeyCreationDateCol, baseEvent.CreationDate()),
				handler.NewCol(AuthNKeyChangeDateCol, baseEvent.CreationDate()),
				handler.NewCol(AuthNKeyResourceOwnerCol, baseEvent.Aggregate().ResourceOwner),
				handler.NewCol(AuthNKeyAggregateIDCol, baseEvent.Aggregate().ID),
				handler.NewCol(AuthNKeySequenceCol, baseEvent.Sequence()),
				handler.NewCol(AuthNKeyObjectIDCol, objectID),
				handler.NewCol(AuthNKeyExpirationCol, expiration),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AuthNKeyPublicIDCol, keyID),
				handler.NewCol(AuthNKeyPublicIdentifierCol, identifier),
				handler.NewCol(AuthNKeyPublicKeyCol, publicKey),
			},
			crdb.WithTableSuffix("public"),
		),
	), nil
}

func (p *AuthNKeyProjection) reduceAuthNKeyRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	var condition handler.Condition
	switch e := event.(type) {
	case *project.ApplicationKeyRemovedEvent:
		condition = handler.NewCond(AuthNKeyIDCol, e.KeyID)
	case *project.APIConfigChangedEvent:
		if e.AuthMethodType != nil && *e.AuthMethodType != domain.APIAuthMethodTypePrivateKeyJWT {
			condition = handler.NewCond(AuthNKeyObjectIDCol, e.AppID)
		}
	case *project.OIDCConfigChangedEvent:
		if e.AuthMethodType != nil && *e.AuthMethodType != domain.OIDCAuthMethodTypePrivateKeyJWT {
			condition = handler.NewCond(AuthNKeyObjectIDCol, e.AppID)
		}
	case *project.ApplicationRemovedEvent:
		condition = handler.NewCond(AuthNKeyObjectIDCol, e.AppID)
	case *project.ProjectRemovedEvent:
		condition = handler.NewCond(AuthNKeyAggregateIDCol, e.Aggregate().ID)
	case *user.MachineKeyRemovedEvent:
		condition = handler.NewCond(AuthNKeyIDCol, e.KeyID)
	case *user.UserRemovedEvent:
		condition = handler.NewCond(AuthNKeyAggregateIDCol, e.Aggregate().ID)
	default:
		logging.LogWithFields("PROJE-Sfdg3", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{project.ApplicationKeyAddedEventType, user.MachineKeyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Dn41s", "reduce.wrong.event.type")
	}
	if condition.Name == "" {
		return crdb.NewNoOpStatement(event), nil
	}
	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{condition},
	), nil
}
