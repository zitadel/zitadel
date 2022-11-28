package projection

import (
	"context"
	"math"

	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	QuotaTable                    = "projections.instance_quotas"
	QuotaNotificationsTableSuffix = "notifications"

	QuotaCreationDateCol              = "creation_date"
	QuotaChangeDateCol                = "change_date"
	QuotaResourceOwnerCol             = "resource_owner"
	QuotaInstanceIDCol                = "instance_id"
	QuotaSequenceCol                  = "sequence"
	QuotaUnitCol                      = "unit"
	QuotaFromCol                      = "from"
	QuotaIntervalCol                  = "interval"
	QuotaAmountCol                    = "amount"
	QuotaLimitationBlockMessageCol    = "limitation_block_message"
	QuotaLimitationBlockHTTPStatusCol = "limitation_block_http_status"
	QuotaLimitationBlockGRPCStatusCol = "limitation_block_grpc_status"
	QuotaLimitationCookieValCol       = "limitation_cookie_val"
	QuotaLimitationRedirectURLCol     = "limitation_redirect_url"
	QuotaUsedAbsoluteCol              = "used_absolute"
	QuotaUsedRelativeCol              = "used_relative"
	QuotaLimitingCol                  = "limiting"

	QuotaNotificationIdCol         = "id"
	QuotaNotificationInstanceIDCol = "instance_id"
	QuotaNotificationUnitCol       = "quota_unit"
	QuotaNotificationCallURLCol    = "call_url"
	QuotaNotificationPercentCol    = "percent"
	QuotaNotificationRepeatCol     = "repeat"
)

type quotaProjection struct {
	crdb.StatementHandler
}

func newQuotaProjection(ctx context.Context, config crdb.StatementHandlerConfig) *quotaProjection {
	p := new(quotaProjection)
	config.ProjectionName = QuotaTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaCreationDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaChangeDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaResourceOwnerCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaInstanceIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaSequenceCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaUnitCol, crdb.ColumnTypeEnum),
				//			crdb.NewColumn(QuotaFromCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaIntervalCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaAmountCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaLimitationBlockMessageCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaLimitationBlockHTTPStatusCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaLimitationBlockGRPCStatusCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaLimitationCookieValCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaLimitationRedirectURLCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaUsedAbsoluteCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaUsedRelativeCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaLimitingCol, crdb.ColumnTypeBool),
			},
			crdb.NewPrimaryKey(QuotaInstanceIDCol, QuotaUnitCol),
			crdb.WithIndex(crdb.NewIndex("quotas_ro_idx", []string{QuotaResourceOwnerCol})),
		),
		crdb.NewSuffixedTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaNotificationIdCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationInstanceIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationUnitCol, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaNotificationCallURLCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationPercentCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaNotificationRepeatCol, crdb.ColumnTypeBool),
			},
			crdb.NewPrimaryKey(QuotaNotificationInstanceIDCol, QuotaNotificationUnitCol, QuotaNotificationIdCol),
			QuotaNotificationsTableSuffix,
			crdb.WithForeignKey(
				crdb.NewForeignKey(
					"fk_instance_quotas_notifications_ref_instance_quotas",
					[]string{QuotaNotificationInstanceIDCol, QuotaNotificationUnitCol},
					[]string{QuotaInstanceIDCol, QuotaUnitCol},
				),
			),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *quotaProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.QuotaAddedEventType,
					Reduce: p.reduceQuotaAdded,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(QuotaInstanceIDCol),
				},
			},
		},
	}
}

func (p *quotaProjection) reduceQuotaAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.QuotaAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dff21", "reduce.wrong.event.type% s", quota.AddedEventType)
	}

	execFuncs := []func(eventstore.Event) crdb.Exec{
		crdb.AddCreateStatement([]handler.Column{
			handler.NewCol(QuotaCreationDateCol, e.CreationDate()),
			handler.NewCol(QuotaChangeDateCol, e.CreationDate()),
			handler.NewCol(QuotaResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(QuotaInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(QuotaSequenceCol, e.Sequence()),
			handler.NewCol(QuotaUnitCol, e.Unit),
			//			handler.NewCol(QuotaFromCol, e.From), TODO: Why is it not working?
			handler.NewCol(QuotaIntervalCol, e.Interval),
			handler.NewCol(QuotaAmountCol, e.Amount),
			handler.NewCol(QuotaLimitationBlockMessageCol, e.Limitations.Block.Message),
			handler.NewCol(QuotaLimitationBlockHTTPStatusCol, e.Limitations.Block.HTTPStatus),
			handler.NewCol(QuotaLimitationBlockGRPCStatusCol, e.Limitations.Block.GRPCStatus),
			handler.NewCol(QuotaLimitationCookieValCol, e.Limitations.CookieValue),
			handler.NewCol(QuotaLimitationRedirectURLCol, e.Limitations.RedirectURL),
			handler.NewCol(QuotaUsedAbsoluteCol, 0),
			handler.NewCol(QuotaUsedRelativeCol, 0),
			handler.NewCol(QuotaLimitingCol, false),
		}),
	}

	for _, notification := range e.Notifications {

		execFuncs = append(execFuncs, crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(QuotaNotificationInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(QuotaNotificationUnitCol, e.Unit),
				handler.NewCol(QuotaNotificationIdCol, notification.ID),
				handler.NewCol(QuotaNotificationPercentCol, notification.Percent),
				handler.NewCol(QuotaNotificationRepeatCol, notification.Repeat),
				handler.NewCol(QuotaNotificationCallURLCol, notification.CallURL),
			},
			crdb.WithTableSuffix(QuotaNotificationsTableSuffix),
		))
	}

	return crdb.NewMultiStatement(e, execFuncs...), nil
}

type nextNotification int64

const (
	unknown nextNotification = -1
	due     nextNotification = -2
	done    nextNotification = -3
)

// TODO: think
func nextNotificationCall(amount, percent, used, currentThreshold int64, repeat bool) int64 {

	if nextNotification(currentThreshold) == due || nextNotification(currentThreshold) == done {
		return currentThreshold
	}

	nextThreshold := int64(math.Floor(float64(amount/100))) * percent
	if used < nextThreshold && currentThreshold != -1 {
		return nextThreshold
	}

	if repeat {
		return nextNotificationCall(amount, percent*2, used, currentThreshold, repeat)
	}
	return -1
}
