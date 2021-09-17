package projection

import (
	"context"
	"fmt"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgDomainProjection struct {
	crdb.StatementHandler
}

const (
	orgDomainProjection = "zitadel.projections.org_domains"
)

func NewOrgDomainProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgDomainProjection {
	p := &OrgDomainProjection{}
	config.ProjectionName = orgDomainProjection
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OrgDomainProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgDomainAddedEventType,
					Reduce: p.reduceDomainAdded,
				},
				{
					Event:  org.OrgDomainVerificationAddedEventType,
					Reduce: p.reduceDomainVerificationAdded,
				},
				{
					Event:  org.OrgDomainVerifiedEventType,
					Reduce: p.reduceDomainVerified,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reducePrimaryDomainSet,
				},
				{
					Event:  org.OrgDomainRemovedEventType,
					Reduce: p.reduceDomainRemoved,
				},
			},
		},
	}
}

const (
	domainCreationDateCol   = "creation_date"
	domainChangeDateCol     = "change_date"
	domainSequenceCol       = "sequence"
	domainDomainCol         = "domain"
	domainOrgIDCol          = "org_id"
	domainIsVerifiedCol     = "is_verified"
	domainIsPrimaryCol      = "is_primary"
	domainValidationTypeCol = "validation_type"
)

func (p *OrgDomainProjection) reduceDomainAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-6fXKf", "seq", event.Sequence(), "expectedType", org.OrgDomainAddedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-DM2DI", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(domainCreationDateCol, e.CreationDate()),
			handler.NewCol(domainChangeDateCol, e.CreationDate()),
			handler.NewCol(domainSequenceCol, e.Sequence()),
			handler.NewCol(domainDomainCol, e.Domain),
			handler.NewCol(domainOrgIDCol, e.Aggregate().ID),
			handler.NewCol(domainIsVerifiedCol, false),
			handler.NewCol(domainIsPrimaryCol, false),
			handler.NewCol(domainValidationTypeCol, domain.OrgDomainValidationTypeUnspecified),
		},
	), nil
}

func (p *OrgDomainProjection) reduceDomainVerificationAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerificationAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-2gGSs", "seq", event.Sequence(), "expectedType", org.OrgDomainVerificationAddedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-EBzyu", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(domainChangeDateCol, e.CreationDate()),
			handler.NewCol(domainSequenceCol, e.Sequence()),
			handler.NewCol(domainValidationTypeCol, e.ValidationType),
		},
		[]handler.Condition{
			handler.NewCond(domainDomainCol, e.Domain),
			handler.NewCond(domainOrgIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *OrgDomainProjection) reduceDomainVerified(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		logging.LogWithFields("PROJE-aeGCA", "seq", event.Sequence(), "expectedType", org.OrgDomainVerifiedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-3Rvkr", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(domainChangeDateCol, e.CreationDate()),
			handler.NewCol(domainSequenceCol, e.Sequence()),
			handler.NewCol(domainIsVerifiedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(domainDomainCol, e.Domain),
			handler.NewCond(domainOrgIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *OrgDomainProjection) reducePrimaryDomainSet(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("PROJE-6YjHo", "seq", event.Sequence(), "expectedType", org.OrgDomainPrimarySetEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-aIuei", "reduce.wrong.event.type")
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(domainChangeDateCol, e.CreationDate()),
				handler.NewCol(domainSequenceCol, e.Sequence()),
				handler.NewCol(domainIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(domainOrgIDCol, e.Aggregate().ID),
				handler.NewCond(domainIsPrimaryCol, true),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(domainChangeDateCol, e.CreationDate()),
				handler.NewCol(domainSequenceCol, e.Sequence()),
				handler.NewCol(domainIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(domainDomainCol, e.Domain),
				handler.NewCond(domainOrgIDCol, e.Aggregate().ID),
			},
		),
	), nil
}

func (p *OrgDomainProjection) reduceDomainRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-dDnps", "seq", event.Sequence(), "expectedType", org.OrgDomainRemovedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-gh1Mx", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(domainDomainCol, e.Domain),
			handler.NewCond(domainOrgIDCol, e.Aggregate().ID),
		},
	), nil
}
