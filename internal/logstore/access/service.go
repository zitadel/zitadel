package access

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/logstore"
)

type Config struct {
	Enabled  bool
	Debounce *logstore.DebouncerConfig
}

type Service struct {
	enabled   bool
	ctx       context.Context
	debouncer *logstore.Debouncer
	dbClient  *sql.DB
}

func NewService(ctx context.Context, cfg *Config, dbClient *sql.DB) *Service {
	svc := &Service{
		enabled:  cfg != nil && cfg.Enabled,
		ctx:      ctx,
		dbClient: dbClient,
	}
	if svc.enabled {
		if cfg.Debounce != nil && cfg.Debounce.MinFrequency > 0 && cfg.Debounce.MaxBulkSize > 0 {
			svc.debouncer = logstore.NewDebouncer(ctx, cfg.Debounce, newStorageBulkSink(dbClient))
		}
	}
	return svc
}

func (s *Service) Handle(ctx context.Context, record *logstore.AccessLogRecord) error {
	if !s.enabled {
		return nil
	}

	if s.debouncer != nil {
		s.debouncer.Add(record)
		return nil
	}

	return storeAccessLogs(ctx, s.dbClient, []any{record})
}

func (s *Service) Enabled() bool {
	return s.enabled
}
