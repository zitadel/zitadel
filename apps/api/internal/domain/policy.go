package domain

type PolicyState int32

const (
	PolicyStateUnspecified PolicyState = iota
	PolicyStateActive
	PolicyStateRemoved

	policyStateCount
)

func (s PolicyState) Exists() bool {
	return s != PolicyStateUnspecified && s != PolicyStateRemoved
}
