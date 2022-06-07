package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
)

type iamProjection struct {
	crdb.StatementHandler
}

const (
	IAMProjectionTable = "zitadel.projections.iam"
)

func newIAMProjection(ctx context.Context, config crdb.StatementHandlerConfig) *iamProjection {
	p := &iamProjection{}
	config.ProjectionName = IAMProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *iamProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.GlobalOrgSetEventType,
					Reduce: p.reduceGlobalOrgSet,
				},
				{
					Event:  iam.ProjectSetEventType,
					Reduce: p.reduceIAMProjectSet,
				},
				{
					Event:  iam.SetupStartedEventType,
					Reduce: p.reduceSetupEvent,
				},
				{
					Event:  iam.SetupDoneEventType,
					Reduce: p.reduceSetupEvent,
				},
			},
		},
	}
}

type IAMColumn string

const (
	IAMColumnID           = "id"
	IAMColumnChangeDate   = "change_date"
	IAMColumnGlobalOrgID  = "global_org_id"
	IAMColumnProjectID    = "iam_project_id"
	IAMColumnSequence     = "sequence"
	IAMColumnSetUpStarted = "setup_started"
	IAMColumnSetUpDone    = "setup_done"
)

func (p *iamProjection) reduceGlobalOrgSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.GlobalOrgSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-3n89fs", "seq", event.Sequence(), "expectedType", iam.GlobalOrgSetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-2n9f2", "reduce.wrong.event.type")
	}
	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(IAMColumnID, e.Aggregate().ID),
			handler.NewCol(IAMColumnChangeDate, e.CreationDate()),
			handler.NewCol(IAMColumnSequence, e.Sequence()),
			handler.NewCol(IAMColumnGlobalOrgID, e.OrgID),
		},
	), nil
}

func (p *iamProjection) reduceIAMProjectSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.ProjectSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-2j9fw", "seq", event.Sequence(), "expectedType", iam.ProjectSetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-30o0e", "reduce.wrong.event.type")
	}
	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(IAMColumnID, e.Aggregate().ID),
			handler.NewCol(IAMColumnChangeDate, e.CreationDate()),
			handler.NewCol(IAMColumnSequence, e.Sequence()),
			handler.NewCol(IAMColumnProjectID, e.ProjectID),
		},
	), nil
}

func (p *iamProjection) reduceSetupEvent(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SetupStepEvent)
	if !ok {
		logging.LogWithFields("HANDL-39fjw", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{iam.SetupDoneEventType, iam.SetupStartedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-d9nfw", "reduce.wrong.event.type")
	}
	columns := []handler.Column{
		handler.NewCol(IAMColumnID, e.Aggregate().ID),
		handler.NewCol(IAMColumnChangeDate, e.CreationDate()),
		handler.NewCol(IAMColumnSequence, e.Sequence()),
	}
	if e.EventType == iam.SetupStartedEventType {
		columns = append(columns, handler.NewCol(IAMColumnSetUpStarted, e.Step))
	} else {
		columns = append(columns, handler.NewCol(IAMColumnSetUpDone, e.Step))
	}
	return crdb.NewUpsertStatement(
		e,
		columns,
	), nil
}
