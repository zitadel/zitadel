package mfa

type State int32

const (
	StateUnspecified State = iota
	StateNotReady
	StateReady

	stateCount
)

func (f State) Valid() bool {
	return f >= 0 && f < stateCount
}
