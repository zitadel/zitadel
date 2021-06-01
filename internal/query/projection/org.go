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

type OrgProjection struct {
	crdb.StatementHandler
}

func NewOrgProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgProjection {
	p := &OrgProjection{}
	config.ProjectionName = "projections.orgs"
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OrgProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: "org",
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgAddedEventType,
					Reduce: p.reduceOrgAdded,
				},
				{
					Event:  org.OrgChangedEventType,
					Reduce: p.reduceOrgChanged,
				},
				{
					Event:  org.OrgDeactivatedEventType,
					Reduce: p.reduceOrgDeactivated,
				},
				{
					Event:  org.OrgReactivatedEventType,
					Reduce: p.reduceOrgReactivated,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reducePrimaryDomainSet,
				},
			},
		},
	}
}

const (
	orgIDCol            = "id"
	orgCreationDateCol  = "creation_date"
	orgChangeDateCol    = "change_date"
	orgResourceOwnerCol = "resource_owner"
	orgStateCol         = "org_state"
	orgSequenceCol      = "sequence"
	orgDomainCol        = "domain"
	orgNameCol          = "name"
)

func (p *OrgProjection) reduceOrgAdded(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", org.OrgAddedEventType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return []handler.Statement{
		crdb.NewCreateStatement([]handler.Column{
			handler.NewCol(orgIDCol, e.Aggregate().ID),
			handler.NewCol(orgCreationDateCol, e.CreationDate()),
			handler.NewCol(orgChangeDateCol, e.CreationDate()),
			handler.NewCol(orgResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(orgSequenceCol, e.Sequence()),
			handler.NewCol(orgNameCol, e.Name),
			handler.NewCol(orgStateCol, domain.OrgStateActive),
		},
			event.Sequence(),
			event.PreviousSequence(),
		),
	}, nil
}

func (p *OrgProjection) reduceOrgChanged(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-q4oq8", "seq", event.Sequence, "expected", org.OrgChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bg8oM", "reduce.wrong.event.type")
	}
	values := []handler.Column{
		handler.NewCol(orgChangeDateCol, e.CreationDate()),
		handler.NewCol(orgSequenceCol, e.Sequence()),
	}
	if e.Name != "" {
		values = append(values, handler.NewCol(orgNameCol, e.Name))
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				handler.NewCol(orgIDCol, e.Aggregate().ID),
			},
			values,
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func (p *OrgProjection) reduceOrgDeactivated(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-1gwdc", "seq", event.Sequence, "expectedType", org.OrgDeactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BApK4", "reduce.wrong.event.type")
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				handler.NewCol(orgIDCol, e.Aggregate().ID),
			},
			[]handler.Column{
				handler.NewCol(orgChangeDateCol, e.CreationDate()),
				handler.NewCol(orgSequenceCol, e.Sequence()),
				handler.NewCol(orgStateCol, domain.OrgStateInactive),
			},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func (p *OrgProjection) reduceOrgReactivated(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Vjwiy", "seq", event.Sequence, "expectedType", org.OrgReactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-o37De", "reduce.wrong.event.type")
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				handler.NewCol(orgIDCol, e.Aggregate().ID),
			},
			[]handler.Column{
				handler.NewCol(orgChangeDateCol, e.CreationDate()),
				handler.NewCol(orgSequenceCol, e.Sequence()),
				handler.NewCol(orgStateCol, domain.OrgStateActive),
			},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}

func (p *OrgProjection) reducePrimaryDomainSet(event eventstore.EventReader) ([]handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("HANDL-79OhB", "seq", event.Sequence, "expectedType", org.OrgDomainPrimarySetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4TbKT", "reduce.wrong.event.type")
	}
	return []handler.Statement{
		crdb.NewUpdateStatement(
			[]handler.Column{
				handler.NewCol(orgIDCol, e.Aggregate().ID),
			},
			[]handler.Column{
				handler.NewCol(orgChangeDateCol, e.CreationDate()),
				handler.NewCol(orgSequenceCol, e.Sequence()),
				handler.NewCol(orgDomainCol, e.Domain),
			},
			e.Sequence(),
			e.PreviousSequence(),
		),
	}, nil
}
