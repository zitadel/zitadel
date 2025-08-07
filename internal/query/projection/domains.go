package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	DomainsTable = "zitadel.domains"

	DomainsIDCol             = "id"
	DomainsInstanceIDCol     = "instance_id"
	DomainsOrgIDCol          = "org_id"
	DomainsDomainCol         = "domain"
	DomainsIsVerifiedCol     = "is_verified"
	DomainsIsPrimaryCol      = "is_primary"
	DomainsValidationTypeCol = "validation_type"
	DomainsCreatedAtCol      = "created_at"
	DomainsUpdatedAtCol      = "updated_at"
	DomainsDeletedAtCol      = "deleted_at"
)

type domainsProjection struct{}

func newDomainsProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(domainsProjection))
}

func (*domainsProjection) Name() string {
	return DomainsTable
}

func (*domainsProjection) Init() *old_handler.Check {
	// The table is created by migration, so we just need to check it exists
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DomainsIDCol, handler.ColumnTypeText),
			handler.NewColumn(DomainsInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(DomainsOrgIDCol, handler.ColumnTypeText),
			handler.NewColumn(DomainsDomainCol, handler.ColumnTypeText),
			handler.NewColumn(DomainsIsVerifiedCol, handler.ColumnTypeBool),
			handler.NewColumn(DomainsIsPrimaryCol, handler.ColumnTypeBool),
			handler.NewColumn(DomainsValidationTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(DomainsCreatedAtCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DomainsUpdatedAtCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DomainsDeletedAtCol, handler.ColumnTypeTimestamp),
		},
			handler.NewPrimaryKey(DomainsIDCol),
		),
	)
}

func (p *domainsProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgDomainAddedEventType,
					Reduce: p.reduceOrgDomainAdded,
				},
				{
					Event:  org.OrgDomainVerificationAddedEventType,
					Reduce: p.reduceOrgDomainVerificationAdded,
				},
				{
					Event:  org.OrgDomainVerifiedEventType,
					Reduce: p.reduceOrgDomainVerified,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reduceOrgPrimaryDomainSet,
				},
				{
					Event:  org.OrgDomainRemovedEventType,
					Reduce: p.reduceOrgDomainRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceDomainAddedEventType,
					Reduce: p.reduceInstanceDomainAdded,
				},
				{
					Event:  instance.InstanceDomainPrimarySetEventType,
					Reduce: p.reduceInstancePrimaryDomainSet,
				},
				{
					Event:  instance.InstanceDomainRemovedEventType,
					Reduce: p.reduceInstanceDomainRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: p.reduceInstanceRemoved,
				},
			},
		},
	}
}

// Organization domain event handlers

func (p *domainsProjection) reduceOrgDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DM2DI", "reduce.wrong.event.type %s", org.OrgDomainAddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(DomainsOrgIDCol, e.Aggregate().ID),
			handler.NewCol(DomainsDomainCol, e.Domain),
			handler.NewCol(DomainsIsVerifiedCol, false),
			handler.NewCol(DomainsIsPrimaryCol, false),
			handler.NewCol(DomainsValidationTypeCol, domain.OrgDomainValidationTypeUnspecified),
			handler.NewCol(DomainsCreatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
		},
	), nil
}

func (p *domainsProjection) reduceOrgDomainVerificationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerificationAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-EBzyu", "reduce.wrong.event.type %s", org.OrgDomainVerificationAddedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsValidationTypeCol, e.ValidationType),
		},
		[]handler.Condition{
			handler.NewCond(DomainsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainsOrgIDCol, e.Aggregate().ID),
			handler.NewCond(DomainsDomainCol, e.Domain),
			handler.NewCond(DomainsDeletedAtCol, nil),
		},
	), nil
}

func (p *domainsProjection) reduceOrgDomainVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-3Rvkr", "reduce.wrong.event.type %s", org.OrgDomainVerifiedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsIsVerifiedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(DomainsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainsOrgIDCol, e.Aggregate().ID),
			handler.NewCond(DomainsDomainCol, e.Domain),
			handler.NewCond(DomainsDeletedAtCol, nil),
		},
	), nil
}

func (p *domainsProjection) reduceOrgPrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-aIuei", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
				handler.NewCol(DomainsIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(DomainsInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(DomainsOrgIDCol, e.Aggregate().ID),
				handler.NewCond(DomainsIsPrimaryCol, true),
				handler.NewCond(DomainsDeletedAtCol, nil),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
				handler.NewCol(DomainsIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(DomainsInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(DomainsOrgIDCol, e.Aggregate().ID),
				handler.NewCond(DomainsDomainCol, e.Domain),
				handler.NewCond(DomainsDeletedAtCol, nil),
			},
		),
	), nil
}

func (p *domainsProjection) reduceOrgDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gh1Mx", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsDeletedAtCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(DomainsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainsOrgIDCol, e.Aggregate().ID),
			handler.NewCond(DomainsDomainCol, e.Domain),
		},
	), nil
}

func (p *domainsProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-dMUKJ", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsDeletedAtCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(DomainsInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainsOrgIDCol, e.Aggregate().ID),
			handler.NewCond(DomainsDeletedAtCol, nil),
		},
	), nil
}

// Instance domain event handlers

func (p *domainsProjection) reduceInstanceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-38nNf", "reduce.wrong.event.type %s", instance.InstanceDomainAddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsInstanceIDCol, e.Aggregate().ID),
			handler.NewCol(DomainsOrgIDCol, nil), // Instance domains have no org_id
			handler.NewCol(DomainsDomainCol, e.Domain),
			handler.NewCol(DomainsIsVerifiedCol, true), // Instance domains are always verified
			handler.NewCol(DomainsIsPrimaryCol, false),
			handler.NewCol(DomainsValidationTypeCol, nil), // Instance domains have no validation type
			handler.NewCol(DomainsCreatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
		},
	), nil
}

func (p *domainsProjection) reduceInstancePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-f8nlw", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
				handler.NewCol(DomainsIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(DomainsInstanceIDCol, e.Aggregate().ID),
				handler.NewCond(DomainsOrgIDCol, nil), // Instance domains
				handler.NewCond(DomainsIsPrimaryCol, true),
				handler.NewCond(DomainsDeletedAtCol, nil),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
				handler.NewCol(DomainsIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(DomainsInstanceIDCol, e.Aggregate().ID),
				handler.NewCond(DomainsOrgIDCol, nil), // Instance domains
				handler.NewCond(DomainsDomainCol, e.Domain),
				handler.NewCond(DomainsDeletedAtCol, nil),
			},
		),
	), nil
}

func (p *domainsProjection) reduceInstanceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-388Nk", "reduce.wrong.event.type %s", instance.InstanceDomainRemovedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsDeletedAtCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(DomainsInstanceIDCol, e.Aggregate().ID),
			handler.NewCond(DomainsOrgIDCol, nil), // Instance domains
			handler.NewCond(DomainsDomainCol, e.Domain),
		},
	), nil
}

func (p *domainsProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-2n9f0", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsDeletedAtCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(DomainsInstanceIDCol, e.Aggregate().ID),
			handler.NewCond(DomainsDeletedAtCol, nil),
		},
	), nil
}