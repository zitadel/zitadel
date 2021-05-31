package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
)

const (
	idColName            = "id"
	creationDateColName  = "creation_date"
	changeDateColName    = "change_date"
	resourceOwnerColName = "resource_owner"
	stateColName         = "org_state"
	sequenceColName      = "sequence"
	domainColName        = "domain"
	nameColName          = "name"
)

var (
	orgReducers = []handler.EventReducer{
		{
			Aggregate: "org",
			Event:     org.OrgAddedEventType,
			Reduce:    orgAddedStmts,
		},
		{
			Aggregate: "org",
			Event:     org.OrgChangedEventType,
			Reduce:    orgChangedStmts,
		},
		{
			Aggregate: "org",
			Event:     org.OrgDeactivatedEventType,
			Reduce:    orgDeactivatedStmts,
		},
		{
			Aggregate: "org",
			Event:     org.OrgReactivatedEventType,
			Reduce:    orgReactivatedStmts,
		},
		{
			Aggregate: "org",
			Event:     org.OrgDomainPrimarySetEventType,
			Reduce:    orgPrimaryDomainStmts,
		},
	}
)

func NewOrgProjection(ctx context.Context, config crdb.StatementHandlerConfig) crdb.StatementHandler {
	config.ProjectionName = "projections.orgs"
	config.Reducers = orgReducers
	return crdb.NewStatementHandler(ctx, config)
}

func orgAddedStmts(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", org.OrgAddedEventType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}

	return []handler.Statement{
		crdb.NewCreateStatement([]handler.Column{
			{
				Name:  idColName,
				Value: e.Aggregate().ID,
			},
			{
				Name:  creationDateColName,
				Value: e.CreationDate(),
			},
			{
				Name:  changeDateColName,
				Value: e.CreationDate(),
			},
			{
				Name:  resourceOwnerColName,
				Value: e.Aggregate().ResourceOwner,
			},
			{
				Name:  sequenceColName,
				Value: e.Sequence(),
			},
			{
				Name:  nameColName,
				Value: e.Name,
			},
			{
				Name:  stateColName,
				Value: domain.OrgStateActive,
			},
		},
			event.Sequence(),
			event.PreviousSequence(),
		),
	}, nil
}

func orgChangedStmts(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-q4oq8", "seq", event.Sequence, "expected", org.OrgChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bg8oM", "reduce.wrong.event.type")
	}
	values := []handler.Column{
		{
			Name:  changeDateColName,
			Value: e.CreationDate(),
		},
		{
			Name:  sequenceColName,
			Value: e.Sequence(),
		},
	}
	if e.Name != "" {
		values = append(values, handler.Column{
			Name:  nameColName,
			Value: e.Name,
		})
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				{
					Name:  idColName,
					Value: e.Aggregate().ID,
				},
			},
			values,
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func orgReactivatedStmts(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Vjwiy", "seq", event.Sequence, "expectedType", org.OrgReactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-o37De", "reduce.wrong.event.type")
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				{
					Name:  idColName,
					Value: e.Aggregate().ID,
				},
			},
			[]handler.Column{
				{
					Name:  changeDateColName,
					Value: e.CreationDate(),
				},
				{
					Name:  sequenceColName,
					Value: e.Sequence(),
				},
				{
					Name:  stateColName,
					Value: domain.OrgStateActive,
				},
			},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func orgDeactivatedStmts(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-1gwdc", "seq", event.Sequence, "expectedType", org.OrgDeactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BApK4", "reduce.wrong.event.type")
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				{
					Name:  idColName,
					Value: e.Aggregate().ID,
				},
			},
			[]handler.Column{
				{
					Name:  changeDateColName,
					Value: e.CreationDate(),
				},
				{
					Name:  sequenceColName,
					Value: e.Sequence(),
				},
				{
					Name:  stateColName,
					Value: domain.OrgStateInactive,
				},
			},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func orgPrimaryDomainStmts(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("HANDL-79OhB", "seq", event.Sequence, "expectedType", org.OrgDomainPrimarySetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4TbKT", "reduce.wrong.event.type")
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				{
					Name:  idColName,
					Value: e.Aggregate().ID,
				},
			},
			[]handler.Column{
				{
					Name:  changeDateColName,
					Value: e.CreationDate(),
				},
				{
					Name:  sequenceColName,
					Value: e.Sequence(),
				},
				{
					Name:  nameColName,
					Value: e.Domain,
				},
			},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}
