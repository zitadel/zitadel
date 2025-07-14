package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/database"
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
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DomainsIDCol, handler.ColumnTypeInt64, handler.Nullable()),
			handler.NewColumn(DomainsInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(DomainsOrgIDCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(DomainsDomainCol, handler.ColumnTypeText),
			handler.NewColumn(DomainsIsVerifiedCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(DomainsIsPrimaryCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(DomainsValidationTypeCol, handler.ColumnTypeInt16, handler.Nullable()),
			handler.NewColumn(DomainsCreatedAtCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DomainsUpdatedAtCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DomainsDeletedAtCol, handler.ColumnTypeTimestamp, handler.Nullable()),
		},
			handler.NewPrimaryKey(DomainsIDCol),
			handler.WithIndex(
				handler.NewIndex("idx_domains_instance_org_domain", 
					[]string{DomainsInstanceIDCol, DomainsOrgIDCol, DomainsDomainCol},
					handler.WithIndexWhere("deleted_at IS NULL"),
				),
			),
			handler.WithIndex(
				handler.NewIndex("idx_domains_instance_id", 
					[]string{DomainsInstanceIDCol},
					handler.WithIndexWhere("deleted_at IS NULL"),
				),
			),
			handler.WithIndex(
				handler.NewIndex("idx_domains_domain", 
					[]string{DomainsDomainCol},
					handler.WithIndexWhere("deleted_at IS NULL"),
				),
			),
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
					Reduce: p.reduceOrgDomainPrimarySet,
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
					Reduce: p.reduceInstanceDomainPrimarySet,
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

// Organization domain event reducers

func (p *domainsProjection) reduceOrgDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D1m8s", "reduce.wrong.event.type %s", org.OrgDomainAddedEventType)
	}

	return p.insertDomain(
		e,
		e.Aggregate().InstanceID,
		&e.Aggregate().ID, // org_id
		e.Domain,
		false, // not verified initially
		false, // not primary initially
		&domain.OrgDomainValidationTypeUnspecified,
	), nil
}

func (p *domainsProjection) reduceOrgDomainVerificationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerificationAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D2m8s", "reduce.wrong.event.type %s", org.OrgDomainVerificationAddedEventType)
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D3m8s", "reduce.wrong.event.type %s", org.OrgDomainVerifiedEventType)
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

func (p *domainsProjection) reduceOrgDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D4m8s", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
	}

	return handler.NewMultiStatement(
		e,
		// First unset all primary domains for this organization
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
		// Then set the new primary domain
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D5m8s", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
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
			handler.NewCond(DomainsDeletedAtCol, nil),
		},
	), nil
}

func (p *domainsProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D6m8s", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
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

// Instance domain event reducers

func (p *domainsProjection) reduceInstanceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D7m8s", "reduce.wrong.event.type %s", instance.InstanceDomainAddedEventType)
	}

	return p.insertDomain(
		e,
		e.Aggregate().ID, // instance_id
		nil,              // org_id is nil for instance domains
		e.Domain,
		true,  // instance domains are always verified
		false, // not primary initially
		nil,   // validation_type is nil for instance domains
	), nil
}

func (p *domainsProjection) reduceInstanceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D8m8s", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}

	return handler.NewMultiStatement(
		e,
		// First unset all primary domains for this instance (where org_id is null)
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
				handler.NewCol(DomainsIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(DomainsInstanceIDCol, e.Aggregate().ID),
				handler.NewCond(DomainsOrgIDCol, nil),
				handler.NewCond(DomainsIsPrimaryCol, true),
				handler.NewCond(DomainsDeletedAtCol, nil),
			},
		),
		// Then set the new primary domain
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
				handler.NewCol(DomainsIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(DomainsInstanceIDCol, e.Aggregate().ID),
				handler.NewCond(DomainsOrgIDCol, nil),
				handler.NewCond(DomainsDomainCol, e.Domain),
				handler.NewCond(DomainsDeletedAtCol, nil),
			},
		),
	), nil
}

func (p *domainsProjection) reduceInstanceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D9m8s", "reduce.wrong.event.type %s", instance.InstanceDomainRemovedEventType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DomainsUpdatedAtCol, e.CreationDate()),
			handler.NewCol(DomainsDeletedAtCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(DomainsInstanceIDCol, e.Aggregate().ID),
			handler.NewCond(DomainsOrgIDCol, nil),
			handler.NewCond(DomainsDomainCol, e.Domain),
			handler.NewCond(DomainsDeletedAtCol, nil),
		},
	), nil
}

func (p *domainsProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJ-D10m8s", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
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

// Helper method to insert a domain
func (p *domainsProjection) insertDomain(
	event eventstore.Event,
	instanceID string,
	orgID *string,
	domain string,
	isVerified bool,
	isPrimary bool,
	validationType *domain.OrgDomainValidationType,
) *handler.Statement {
	columns := []handler.Column{
		handler.NewCol(DomainsInstanceIDCol, instanceID),
		handler.NewCol(DomainsDomainCol, domain),
		handler.NewCol(DomainsIsVerifiedCol, isVerified),
		handler.NewCol(DomainsIsPrimaryCol, isPrimary),
		handler.NewCol(DomainsCreatedAtCol, event.CreationDate()),
		handler.NewCol(DomainsUpdatedAtCol, event.CreationDate()),
	}

	if orgID != nil {
		columns = append(columns, handler.NewCol(DomainsOrgIDCol, *orgID))
	}

	if validationType != nil {
		columns = append(columns, handler.NewCol(DomainsValidationTypeCol, *validationType))
	}

	return handler.NewCreateStatement(event, columns)
}