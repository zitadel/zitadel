package debouncer

import (
	"time"
)

type Config struct {
	MinFrequency time.Duration
	MaxBulkSize  uint
}

type Shipper interface {
	Ship([]any)
}

type ShipFunc func([]any)

func (s ShipFunc) Ship(items []any) {
	s.Ship(items)
}

type debouncer interface {
	Add(item any)
}

type Service struct {
	debouncer
}

// New returns a debouncer service, that caches items added until cfg.MinFrequency or cfg.MaxBulkSize is reached.
// Then, the ship function is called with all items added, then the items, ticker and count are reset.
// If cfg.MinFrequency is 0 or cfg.MaxBulkSize is 0, the items are shipped synchronously.
// If cfg is nil, nothing is done with items sent to Add()
func New(cfg *Config, ship Shipper) *Service {

	if cfg == nil {
		return &Service{debouncer: newNoopDebouncer()}
	}

	// TODO: I think MaxBulkSize = 0 should mean math.MaxInt bulk size
	if cfg.MinFrequency == 0 || cfg.MaxBulkSize == 0 {
		return &Service{debouncer: newSyncDebouncer(ship)}
	}

	return &Service{debouncer: newAsyncDebouncer(cfg, ship)}
}
