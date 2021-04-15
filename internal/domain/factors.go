package domain

type SecondFactorType int32

const (
	SecondFactorTypeUnspecified SecondFactorType = iota
	SecondFactorTypeOTP
	SecondFactorTypeU2F

	secondFactorCount
)

func SecondFactorTypes() []SecondFactorType {
	types := make([]SecondFactorType, 0, secondFactorCount-1)
	for i := SecondFactorTypeUnspecified + 1; i < secondFactorCount; i++ {
		types = append(types, i)
	}
	return types
}

type MultiFactorType int32

const (
	MultiFactorTypeUnspecified MultiFactorType = iota
	MultiFactorTypeU2FWithPIN

	multiFactorCount
)

func MultiFactorTypes() []MultiFactorType {
	types := make([]MultiFactorType, 0, multiFactorCount-1)
	for i := MultiFactorTypeUnspecified + 1; i < multiFactorCount; i++ {
		types = append(types, i)
	}
	return types
}

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
