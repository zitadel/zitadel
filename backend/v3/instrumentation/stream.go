package instrumentation

import "sync"

type Stream int

//go:generate enumer -type=Stream -trimprefix=Stream -transform=snake -text
const (
	StreamRuntime      Stream = iota // Top-level commands, such as starting the application or running migrations.
	StreamReady                      // Readiness and liveness checks.
	StreamRequest                    // API request handling.
	StreamEventPusher                // Event pushing to the database.
	StreamEventHandler               // Event handling and processing.
	StreamQueue                      // Queue operations and job processing.
)

var enabledStreams sync.Map

func EnableStreams(streams ...Stream) {
	enabledStreams.Clear()
	for _, stream := range streams {
		enabledStreams.Store(stream, struct{}{})
	}
}

func IsStreamEnabled(stream Stream) bool {
	_, ok := enabledStreams.Load(stream)
	return ok
}
