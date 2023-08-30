package logstore

import (
	"context"
	"math"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const handleThresholdTimeout = time.Minute

type UsageStorer[T LogRecord[T]] interface {
	LogEmitter[T]
	QuotaUnit() quota.Unit
}

type Service[T LogRecord[T]] struct {
	commands         Commands
	queries          Queries
	usageStorer      UsageStorer[T]
	enabledSinks     []*emitter[T]
	sinkEnabled      bool
	reportingEnabled bool
}

type Queries interface {
	GetQuota(ctx context.Context, instanceID string, unit quota.Unit) (qu *query.Quota, err error)
	GetQuotaUsage(ctx context.Context, instanceID string, unit quota.Unit, periodStart time.Time) (usage uint64, err error)
	GetDueQuotaNotifications(ctx context.Context, instanceID string, unit quota.Unit, qu *query.Quota, periodStart time.Time, usedAbs uint64) (dueNotifications []*quota.NotificationDueEvent, err error)
}

type Commands interface {
	ReportQuotaUsage(ctx context.Context, dueNotifications []*quota.NotificationDueEvent) error
}

func New[T LogRecord[T]](queries Queries, commands Commands, usageQuerierSink *emitter[T], additionalSink ...*emitter[T]) *Service[T] {
	var usageStorer UsageStorer[T]
	if usageQuerierSink != nil {
		usageStorer = usageQuerierSink.emitter.(UsageStorer[T])
	}
	svc := &Service[T]{
		commands:         commands,
		queries:          queries,
		reportingEnabled: usageQuerierSink != nil && usageQuerierSink.enabled,
		usageStorer:      usageStorer,
	}
	for _, s := range append([]*emitter[T]{usageQuerierSink}, additionalSink...) {
		if s != nil && s.enabled {
			svc.enabledSinks = append(svc.enabledSinks, s)
		}
	}
	svc.sinkEnabled = len(svc.enabledSinks) > 0
	return svc
}

func (s *Service[T]) Enabled() bool {
	return s.sinkEnabled
}

func (s *Service[T]) Handle(ctx context.Context, record T) {
	for _, sink := range s.enabledSinks {
		logging.OnError(sink.Emit(ctx, record.Normalize())).WithField("record", record).Warn("failed to emit log record")
	}
}

func (s *Service[T]) Limit(ctx context.Context, instanceID string) *uint64 {
	var err error
	defer func() {
		logging.OnError(err).Warn("failed to check if usage should be limited")
	}()
	if !s.reportingEnabled || instanceID == "" {
		return nil
	}
	quotaUnit := s.usageStorer.QuotaUnit()
	q, err := s.queries.GetQuota(ctx, instanceID, quotaUnit)
	if caos_errors.IsNotFound(err) {
		err = nil
		return nil
	}
	if err != nil {
		return nil
	}
	usage, err := s.queries.GetQuotaUsage(ctx, instanceID, s.usageStorer.QuotaUnit(), q.CurrentPeriodStart)
	if err != nil {
		return nil
	}
	go s.handleThresholds(ctx, instanceID, quotaUnit, q, usage)
	var remaining *uint64
	if q.Limit {
		r := uint64(math.Max(0, float64(q.Amount)-float64(usage)))
		remaining = &r
	}
	return remaining
}

func (s *Service[T]) handleThresholds(ctx context.Context, instanceID string, unit quota.Unit, q *query.Quota, usage uint64) {
	var err error
	defer func() {
		logging.OnError(err).Warn("handling quota thresholds failed")
	}()
	detatchedCtx, cancel := context.WithTimeout(authz.Detach(ctx), handleThresholdTimeout)
	defer cancel()
	notifications, err := s.queries.GetDueQuotaNotifications(detatchedCtx, instanceID, unit, q, q.CurrentPeriodStart, usage)
	if err != nil || len(notifications) == 0 {
		return
	}
	err = s.commands.ReportQuotaUsage(detatchedCtx, notifications)
}
