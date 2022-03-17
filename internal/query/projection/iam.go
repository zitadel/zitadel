package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
)

type IAMProjection struct {
	crdb.StatementHandler
}

const (
	IAMProjectionTable = "zitadel.projections.iam"
)

func NewIAMProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IAMProjection {
	p := new(IAMProjection)
	config.ProjectionName = IAMProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *IAMProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.GlobalOrgSetEventType,
					Reduce: p.reduceGlobalOrgSet,
				},
				{
					Event:  instance.ProjectSetEventType,
					Reduce: p.reduceIAMProjectSet,
				},
				{
					Event:  instance.DefaultLanguageSetEventType,
					Reduce: p.reduceDefaultLanguageSet,
				},
				{
					Event:  instance.SetupStartedEventType,
					Reduce: p.reduceSetupEvent,
				},
				{
					Event:  instance.SetupDoneEventType,
					Reduce: p.reduceSetupEvent,
				},
			},
		},
	}
}

type IAMColumn string

const (
	IAMColumnID              = "id"
	IAMColumnChangeDate      = "change_date"
	IAMColumnGlobalOrgID     = "global_org_id"
	IAMColumnProjectID       = "iam_project_id"
	IAMColumnSequence        = "sequence"
	IAMColumnSetUpStarted    = "setup_started"
	IAMColumnSetUpDone       = "setup_done"
	IAMColumnDefaultLanguage = "default_language"
)

func (p *IAMProjection) reduceGlobalOrgSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.GlobalOrgSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-3n89fs", "seq", event.Sequence(), "expectedType", instance.GlobalOrgSetEventType).Error("wrong event type")
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

func (p *IAMProjection) reduceIAMProjectSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.ProjectSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-2j9fw", "seq", event.Sequence(), "expectedType", instance.ProjectSetEventType).Error("wrong event type")
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

func (p *IAMProjection) reduceDefaultLanguageSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DefaultLanguageSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-3n9le", "seq", event.Sequence(), "expectedType", instance.DefaultLanguageSetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-30o0e", "reduce.wrong.event.type")
	}
	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(IAMColumnID, e.Aggregate().ID),
			handler.NewCol(IAMColumnChangeDate, e.CreationDate()),
			handler.NewCol(IAMColumnSequence, e.Sequence()),
			handler.NewCol(IAMColumnDefaultLanguage, e.Language.String()),
		},
	), nil
}

func (p *IAMProjection) reduceSetupEvent(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SetupStepEvent)
	if !ok {
		logging.LogWithFields("HANDL-39fjw", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{instance.SetupDoneEventType, instance.SetupStartedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-d9nfw", "reduce.wrong.event.type")
	}
	columns := []handler.Column{
		handler.NewCol(IAMColumnID, e.Aggregate().ID),
		handler.NewCol(IAMColumnChangeDate, e.CreationDate()),
		handler.NewCol(IAMColumnSequence, e.Sequence()),
	}
	if e.EventType == instance.SetupStartedEventType {
		columns = append(columns, handler.NewCol(IAMColumnSetUpStarted, e.Step))
	} else {
		columns = append(columns, handler.NewCol(IAMColumnSetUpDone, e.Step))
	}
	return crdb.NewUpsertStatement(
		e,
		columns,
	), nil
}
