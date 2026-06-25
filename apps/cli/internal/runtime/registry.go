package runtime

import (
	"sync"
)

var (
	mu    sync.Mutex
	specs []CommandSpec
)

// Register adds a command spec to the global registry.
func Register(s CommandSpec) {
	mu.Lock()
	defer mu.Unlock()
	specs = append(specs, s)
}

// RegisterAll adds multiple command specs to the global registry.
func RegisterAll(ss []CommandSpec) {
	mu.Lock()
	defer mu.Unlock()
	specs = append(specs, ss...)
}

// AllSpecs returns all registered command specs.
func AllSpecs() []CommandSpec {
	mu.Lock()
	defer mu.Unlock()
	out := make([]CommandSpec, len(specs))
	copy(out, specs)
	return out
}
