package access

import (
	"context"
	"database/sql"

	caos_errors "github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/logstore"
)

type Config struct {
	Enabled  bool
	Debounce *logstore.DebouncerConfig
}

type Service struct {
	enabled       bool
	ctx           context.Context
	debouncer     *logstore.Debouncer
	dbClient      *sql.DB
	storedHandler logstore.StoredAccessLogsReducer
	limit         map[string]bool // TODO: dangerous?
}

func NewService(ctx context.Context, cfg *Config, dbClient *sql.DB) *Service {
	svc := &Service{
		enabled:  cfg != nil && cfg.Enabled,
		ctx:      ctx,
		dbClient: dbClient,
		limit:    make(map[string]bool),
	}

	svc.storedHandler = logstore.StoredAccessLogsReducerFunc(func(ctx context.Context, records []*logstore.AccessLogRecord) {
		inst := authenticatedRequests(records)
		for k, v := range inst {
			limit, err := projection.UpdateInstanceQuotaUsage(ctx, dbClient, k, quota.RequestsAllAuthenticated, v)
			if err != nil {
				logging.Warn("updating instance quota usage failed: %s", err.Error())
			}
			if limit {
				svc.limit[k] = true
			}
		}
	})

	if svc.enabled {
		if cfg.Debounce != nil && cfg.Debounce.MinFrequency > 0 && cfg.Debounce.MaxBulkSize > 0 {
			svc.debouncer = logstore.NewDebouncer(ctx, cfg.Debounce, newStorageBulkSink(dbClient, svc.storedHandler))
		}
	}
	return svc
}

func (s *Service) Handle(ctx context.Context, record *logstore.AccessLogRecord) (err error) {
	if !s.enabled {
		return nil
	}

	if record.IsAuthenticated() && s.limit[record.InstanceID] {
		return caos_errors.ThrowError(nil, "QUOTA-9D153", "Errors.Quota.LimitExceeded")
	}

	if s.debouncer != nil {
		s.debouncer.Add(record)
		return nil
	}

	return storeAccessLogs(ctx, s.dbClient, []any{record}, s.storedHandler)
}

func (s *Service) Enabled() bool {
	return s.enabled
}

func authenticatedRequests(records []*logstore.AccessLogRecord) map[string]int64 {
	// TODO: dangerous?
	usage := make(map[string]int64)
	for idx := range records {
		record := records[idx]
		if record.IsAuthenticated() {
			usage[record.InstanceID]++
		}
	}
	return usage
}
