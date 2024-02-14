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
	OrgDomainTable = "projections.org_domains2"

	OrgDomainOrgIDCol          = "org_id"
	OrgDomainInstanceIDCol     = "instance_id"
	OrgDomainCreationDateCol   = "creation_date"
	OrgDomainChangeDateCol     = "change_date"
	OrgDomainSequenceCol       = "sequence"
	OrgDomainDomainCol         = "domain"
	OrgDomainIsVerifiedCol     = "is_verified"
	OrgDomainIsPrimaryCol      = "is_primary"
	OrgDomainValidationTypeCol = "validation_type"
	OrgDomainOwnerRemovedCol   = "owner_removed"
)

type orgDomainProjection struct{}

func newOrgDomainProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(orgDomainProjection))
}

func (*orgDomainProjection) Name() string {
	return OrgDomainTable
}

func (*orgDomainProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(OrgDomainOrgIDCol, handler.ColumnTypeText),
			handler.NewColumn(OrgDomainInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(OrgDomainCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrgDomainChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrgDomainSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(OrgDomainDomainCol, handler.ColumnTypeText),
			handler.NewColumn(OrgDomainIsVerifiedCol, handler.ColumnTypeBool),
			handler.NewColumn(OrgDomainIsPrimaryCol, handler.ColumnTypeBool),
			handler.NewColumn(OrgDomainValidationTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(OrgDomainOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(OrgDomainOrgIDCol, OrgDomainDomainCol, OrgDomainInstanceIDCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{OrgDomainOwnerRemovedCol})),
		),
	)
}

func (p *orgDomainProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Reduce: reduceInstanceRemovedHelper(OrgDomainInstanceIDCol),
				},
			},
		},
	}
}

func (p *orgDomainProjection) reduceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DM2DI", "reduce.wrong.event.type %s", org.OrgDomainAddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgDomainCreationDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
			handler.NewCol(OrgDomainDomainCol, e.Domain),
			handler.NewCol(OrgDomainOrgIDCol, e.Aggregate().ID),
			handler.NewCol(OrgDomainInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(OrgDomainIsVerifiedCol, false),
			handler.NewCol(OrgDomainIsPrimaryCol, false),
			handler.NewCol(OrgDomainValidationTypeCol, domain.OrgDomainValidationTypeUnspecified),
		},
	), nil
}

func (p *orgDomainProjection) reduceDomainVerificationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerificationAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-EBzyu", "reduce.wrong.event.type %s", org.OrgDomainVerificationAddedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
			handler.NewCol(OrgDomainValidationTypeCol, e.ValidationType),
		},
		[]handler.Condition{
			handler.NewCond(OrgDomainDomainCol, e.Domain),
			handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
			handler.NewCond(OrgDomainInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgDomainProjection) reduceDomainVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-3Rvkr", "reduce.wrong.event.type %s", org.OrgDomainVerifiedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
			handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
			handler.NewCol(OrgDomainIsVerifiedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(OrgDomainDomainCol, e.Domain),
			handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
			handler.NewCond(OrgDomainInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgDomainProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-aIuei", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
				handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
				handler.NewCol(OrgDomainIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
				handler.NewCond(OrgDomainIsPrimaryCol, true),
				handler.NewCond(OrgDomainInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(OrgDomainChangeDateCol, e.CreationDate()),
				handler.NewCol(OrgDomainSequenceCol, e.Sequence()),
				handler.NewCol(OrgDomainIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(OrgDomainDomainCol, e.Domain),
				handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
				handler.NewCond(OrgDomainInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *orgDomainProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gh1Mx", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgDomainDomainCol, e.Domain),
			handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
			handler.NewCond(OrgDomainInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgDomainProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-dMUKJ", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgDomainInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(OrgDomainOrgIDCol, e.Aggregate().ID),
		},
	), nil
}
