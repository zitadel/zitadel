package debouncer

type syncDebouncer struct {
	shipper Shipper
}

func newSyncDebouncer(ship Shipper) *syncDebouncer {
	return &syncDebouncer{shipper: ship}
}

func (s *syncDebouncer) Add(item any) {
	sl := []any{item}
	s.shipper.Ship(sl)
}
