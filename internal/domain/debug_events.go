package domain

type DebugEventsState int

const (
	DebugEventsStateUnspecified DebugEventsState = iota
	DebugEventsStateInitial
	DebugEventsStateChanged
	DebugEventsStateRemoved
)

func (state DebugEventsState) Exists() bool {
	return state == DebugEventsStateInitial || state == DebugEventsStateChanged
}
