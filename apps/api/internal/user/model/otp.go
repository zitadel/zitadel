package model

type MFAState int32

const (
	MFAStateUnspecified MFAState = iota
	MFAStateNotReady
	MFAStateReady
)

type MultiFactor struct {
	Type      MFAType
	State     MFAState
	Attribute string
	ID        string
}

type MFAType int32

const (
	MFATypeUnspecified MFAType = iota
	MFATypeOTP
	MFATypeU2F
)
