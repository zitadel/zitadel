package domain

type SecondFactorType int32

const (
	SecondFactorTypeUnspecified SecondFactorType = iota
	SecondFactorTypeOTP
	SecondFactorTypeU2F

	secondFactorCount
)

type MultiFactorType int32

const (
	MultiFactorTypeUnspecified MultiFactorType = iota
	MultiFactorTypeU2FWithPIN

	multiFactorCount
)

type FactorState int32

const (
	FactorStateUnspecified FactorState = iota
	FactorStateActive
	FactorStateRemoved

	factorStateCount
)

func (f SecondFactorType) Valid() bool {
	return f > 0 && f < secondFactorCount
}

func (f MultiFactorType) Valid() bool {
	return f > 0 && f < multiFactorCount
}

func (f FactorState) Valid() bool {
	return f >= 0 && f < factorStateCount
}
