package cachekey

type AuthnKeyIndex int

const (
	InstanceID AuthnKeyIndex = iota
	UserType
	KeyID
)

func (i AuthnKeyIndex) Key() string {
	switch i {
	case InstanceID:
		return "instance_id"
	case UserType:
		return "user_type"
	case KeyID:
		return "key_id"
	default:
		return "unknown"
	}
}
