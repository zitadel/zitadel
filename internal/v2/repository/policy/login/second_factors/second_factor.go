package second_factors

type SecondFactorType int32

const (
	SecondFactorTypeUnspecified SecondFactorType = iota
	SecondFactorTypeOTP
	SecondFactorTypeU2F
)
