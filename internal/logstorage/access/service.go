package access

import (
	"database/sql"

	"github.com/zitadel/zitadel/internal/logstorage/debouncer"
)

type Service struct {
	dbClient        *sql.DB
	debounceService *debouncer.Service
}

func NewAccessLogsStorageService(dbClient *sql.DB) *Service {
	return &Service{
		dbClient:        dbClient,
		debounceService: debounceService,
	}
}

func (s *Service) Handle(entry interface{}) {
	if s.cfg == nil || !s.cfg.Enabled {
		return
	}
}
