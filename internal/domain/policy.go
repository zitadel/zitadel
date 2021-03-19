package domain

type PolicyState int32

const (
	PolicyStateUnspecified PolicyState = iota
	PolicyStateActive
	PolicyStateRemoved

	policyStateCount
)

func (f PolicyState) Valid() bool {
	return f >= 0 && f < policyStateCount
}

func (s PolicyState) Exists() bool {
	return s != PolicyStateUnspecified && s != PolicyStateRemoved
}
