package logstore

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/repository/quota"
)

type UsageStorer[T LogRecord[T]] interface {
	LogEmitter[T]
	QuotaUnit() quota.Unit
}

type Service[T LogRecord[T]] struct {
	queries          Queries
	usageStorer      UsageStorer[T]
	enabledSinks     []*emitter[T]
	sinkEnabled      bool
	reportingEnabled bool
}

type Queries interface {
	GetRemainingQuotaUsage(ctx context.Context, instanceID string, unit quota.Unit) (remaining *uint64, err error)
}

func New[T LogRecord[T]](queries Queries, usageQuerierSink *emitter[T], additionalSink ...*emitter[T]) *Service[T] {
	var usageStorer UsageStorer[T]
	if usageQuerierSink != nil {
		usageStorer = usageQuerierSink.emitter.(UsageStorer[T])
	}
	svc := &Service[T]{
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
	remaining, err := s.queries.GetRemainingQuotaUsage(ctx, instanceID, s.usageStorer.QuotaUnit())
	if err != nil {
		// TODO: shouldn't we just limit then or return the error and decide there?
		return nil
	}
	return remaining
}
