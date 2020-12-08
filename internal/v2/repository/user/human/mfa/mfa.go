package mfa

type MFAState int32

const (
	MFAStateUnspecified MFAState = iota
	MFAStateNotReady
	MFAStateReady

	mfaStateCount
)

func (f MFAState) Valid() bool {
	return f >= 0 && f < mfaStateCount
}
