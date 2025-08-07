package projection

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	HostedLoginTranslationTable = "projections.hosted_login_translations"

	HostedLoginTranslationInstanceIDCol    = "instance_id"
	HostedLoginTranslationCreationDateCol  = "creation_date"
	HostedLoginTranslationChangeDateCol    = "change_date"
	HostedLoginTranslationAggregateIDCol   = "aggregate_id"
	HostedLoginTranslationAggregateTypeCol = "aggregate_type"
	HostedLoginTranslationSequenceCol      = "sequence"
	HostedLoginTranslationLocaleCol        = "locale"
	HostedLoginTranslationFileCol          = "file"
	HostedLoginTranslationEtagCol          = "etag"
)

type hostedLoginTranslationProjection struct{}

func newHostedLoginTranslationProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(hostedLoginTranslationProjection))
}

// Init implements [handler.initializer]
func (p *hostedLoginTranslationProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(HostedLoginTranslationInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(HostedLoginTranslationCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(HostedLoginTranslationChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(HostedLoginTranslationAggregateIDCol, handler.ColumnTypeText),
			handler.NewColumn(HostedLoginTranslationAggregateTypeCol, handler.ColumnTypeText),
			handler.NewColumn(HostedLoginTranslationSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(HostedLoginTranslationLocaleCol, handler.ColumnTypeText),
			handler.NewColumn(HostedLoginTranslationFileCol, handler.ColumnTypeJSONB),
			handler.NewColumn(HostedLoginTranslationEtagCol, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(
				HostedLoginTranslationInstanceIDCol,
				HostedLoginTranslationAggregateIDCol,
				HostedLoginTranslationAggregateTypeCol,
				HostedLoginTranslationLocaleCol,
			),
		),
	)
}

func (hltp *hostedLoginTranslationProjection) Name() string {
	return HostedLoginTranslationTable
}

func (hltp *hostedLoginTranslationProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.HostedLoginTranslationSet,
					Reduce: hltp.reduceSet,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.HostedLoginTranslationSet,
					Reduce: hltp.reduceSet,
				},
			},
		},
	}
}

func (hltp *hostedLoginTranslationProjection) reduceSet(e eventstore.Event) (*handler.Statement, error) {

	switch e := e.(type) {
	case *org.HostedLoginTranslationSetEvent:
		orgEvent := *e
		return handler.NewUpsertStatement(
			&orgEvent,
			[]handler.Column{
				handler.NewCol(HostedLoginTranslationInstanceIDCol, nil),
				handler.NewCol(HostedLoginTranslationAggregateIDCol, nil),
				handler.NewCol(HostedLoginTranslationAggregateTypeCol, nil),
				handler.NewCol(HostedLoginTranslationLocaleCol, nil),
			},
			[]handler.Column{
				handler.NewCol(HostedLoginTranslationInstanceIDCol, orgEvent.Aggregate().InstanceID),
				handler.NewCol(HostedLoginTranslationAggregateIDCol, orgEvent.Aggregate().ID),
				handler.NewCol(HostedLoginTranslationAggregateTypeCol, orgEvent.Aggregate().Type),
				handler.NewCol(HostedLoginTranslationCreationDateCol, handler.OnlySetValueOnInsert(HostedLoginTranslationTable, orgEvent.CreationDate())),
				handler.NewCol(HostedLoginTranslationChangeDateCol, orgEvent.CreationDate()),
				handler.NewCol(HostedLoginTranslationSequenceCol, orgEvent.Sequence()),
				handler.NewCol(HostedLoginTranslationLocaleCol, orgEvent.Language),
				handler.NewCol(HostedLoginTranslationFileCol, orgEvent.Translation),
				handler.NewCol(HostedLoginTranslationEtagCol, hltp.computeEtag(orgEvent.Translation)),
			},
		), nil
	case *instance.HostedLoginTranslationSetEvent:
		instanceEvent := *e
		return handler.NewUpsertStatement(
			&instanceEvent,
			[]handler.Column{
				handler.NewCol(HostedLoginTranslationInstanceIDCol, nil),
				handler.NewCol(HostedLoginTranslationAggregateIDCol, nil),
				handler.NewCol(HostedLoginTranslationAggregateTypeCol, nil),
				handler.NewCol(HostedLoginTranslationLocaleCol, nil),
			},
			[]handler.Column{
				handler.NewCol(HostedLoginTranslationInstanceIDCol, instanceEvent.Aggregate().InstanceID),
				handler.NewCol(HostedLoginTranslationAggregateIDCol, instanceEvent.Aggregate().ID),
				handler.NewCol(HostedLoginTranslationAggregateTypeCol, instanceEvent.Aggregate().Type),
				handler.NewCol(HostedLoginTranslationCreationDateCol, handler.OnlySetValueOnInsert(HostedLoginTranslationTable, instanceEvent.CreationDate())),
				handler.NewCol(HostedLoginTranslationChangeDateCol, instanceEvent.CreationDate()),
				handler.NewCol(HostedLoginTranslationSequenceCol, instanceEvent.Sequence()),
				handler.NewCol(HostedLoginTranslationLocaleCol, instanceEvent.Language),
				handler.NewCol(HostedLoginTranslationFileCol, instanceEvent.Translation),
				handler.NewCol(HostedLoginTranslationEtagCol, hltp.computeEtag(instanceEvent.Translation)),
			},
		), nil
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-AZshaa", "reduce.wrong.event.type %v", []eventstore.EventType{org.HostedLoginTranslationSet})
	}

}

func (hltp *hostedLoginTranslationProjection) computeEtag(translation map[string]any) string {
	hash := md5.Sum(fmt.Append(nil, translation))
	return hex.EncodeToString(hash[:])
}
