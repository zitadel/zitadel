package domain

type SecondFactorType int32

const (
	SecondFactorTypeUnspecified SecondFactorType = iota
	SecondFactorTypeOTP
	SecondFactorTypeU2F
)

type MultiFactorType int32

const (
	MultiFactorTypeUnspecified MultiFactorType = iota
	MultiFactorTypeU2FWithPIN
)

type FactorState int32

const (
	FactorStateUnspecified FactorState = iota
	FactorStateActive
	FactorStateRemoved

	factorStateCount
)

func (f FactorState) Valid() bool {
	return f >= 0 && f < factorStateCount
}
