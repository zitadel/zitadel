package domain

const (
	IAMID = "IAM"
)

type InstanceState int32

const (
	InstanceStateUnspecified InstanceState = iota
	InstanceStateActive
	InstanceStateRemoved

	instanceStateCount
)

func (f InstanceState) Valid() bool {
	return f >= 0 && f < instanceStateCount
}
