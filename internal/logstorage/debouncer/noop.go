package debouncer

type noopDebouncer struct{}

func newNoopDebouncer() *noopDebouncer {
	return &noopDebouncer{}
}

func (*noopDebouncer) Add(_ any) {}
