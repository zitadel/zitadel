package domain

type Gender int32

const (
	GenderUnspecified Gender = iota
	GenderFemale
	GenderMale
	GenderDiverse

	genderCount
)

func (f Gender) Valid() bool {
	return f >= 0 && f < genderCount
}
