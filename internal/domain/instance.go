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

func (s InstanceState) Valid() bool {
	return s >= 0 && s < instanceStateCount
}

func (s InstanceState) Exists() bool {
	return s != InstanceStateUnspecified && s != InstanceStateRemoved
}
