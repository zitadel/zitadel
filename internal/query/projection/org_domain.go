package projection

import (
	"context"
	"fmt"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type orgDomainProjection struct {
	crdb.StatementHandler
}

const (
	OrgDomainTable = "zitadel.projections.org_domains"
)

func newOrgDomainProjection(ctx context.Context, config crdb.StatementHandlerConfig) *orgDomainProjection {
	p := &orgDomainProjection{}
	config.ProjectionName = OrgDomainTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *orgDomainProjection) reducers() []handler.AggregateReducer {
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
	OrgDomainCreationDateCol   = "creation_date"
	OrgDomainChangeDateCol     = "change_date"
	OrgDomainSequenceCol       = "sequence"
	OrgDomainDomainCol         = "domain"
	OrgDomainOrgIDCol          = "org_id"
	OrgDomainIsVerifiedCol     = "is_verified"
	OrgDomainIsPrimaryCol      = "is_primary"
	OrgDomainValidationTypeCol = "validation_type"
)

func (p *orgDomainProjection) reduceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-6fXKf", "seq", event.Sequence(), "expectedType", org.OrgDomainAddedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-DM2DI", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgDomainCreationDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
			handler.NewCol(OrgDomainDomainCol, e.Domain),
			handler.NewCol(OrgDomainOrgIDCol, e.Aggregate().ID),
			handler.NewCol(OrgDomainIsVerifiedCol, false),
			handler.NewCol(OrgDomainIsPrimaryCol, false),
			handler.NewCol(OrgDomainValidationTypeCol, domain.OrgDomainValidationTypeUnspecified),
		},
	), nil
}

func (p *orgDomainProjection) reduceDomainVerificationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerificationAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-2gGSs", "seq", event.Sequence(), "expectedType", org.OrgDomainVerificationAddedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-EBzyu", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
			handler.NewCol(OrgDomainValidationTypeCol, e.ValidationType),
		},
		[]handler.Condition{
			handler.NewCond(OrgDomainDomainCol, e.Domain),
			handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *orgDomainProjection) reduceDomainVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		logging.LogWithFields("PROJE-aeGCA", "seq", event.Sequence(), "expectedType", org.OrgDomainVerifiedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-3Rvkr", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
			handler.NewCol(OrgDomainIsVerifiedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(OrgDomainDomainCol, e.Domain),
			handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *orgDomainProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("PROJE-6YjHo", "seq", event.Sequence(), "expectedType", org.OrgDomainPrimarySetEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-aIuei", "reduce.wrong.event.type")
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
				handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
				handler.NewCol(OrgDomainIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
				handler.NewCond(OrgDomainIsPrimaryCol, true),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
				handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
				handler.NewCol(OrgDomainIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(OrgDomainDomainCol, e.Domain),
				handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
			},
		),
	), nil
}

func (p *orgDomainProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-dDnps", "seq", event.Sequence(), "expectedType", org.OrgDomainRemovedEventType, "gottenType", fmt.Sprintf("%T", event)).Error("unexpected event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-gh1Mx", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgDomainDomainCol, e.Domain),
			handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
		},
	), nil
}
