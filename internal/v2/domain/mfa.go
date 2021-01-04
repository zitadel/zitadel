package domain

type MFAState int32

const (
	MFAStateUnspecified MFAState = iota
	MFAStateNotReady
	MFAStateReady

	stateCount
)

func (f MFAState) Valid() bool {
	return f >= 0 && f < stateCount
}
