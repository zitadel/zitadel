package deviceauth

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	UniqueUserCode      = "user_code"
	UniqueDeviceCode    = "device_code"
	DuplicateUserCode   = "Errors.DeviceUserCode.AlreadyExists"
	DuplicateDeviceCode = "Errors.DeviceCode.AlreadyExists"
)

func NewAddUniqueConstraints(deviceCode, userCode string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(
			UniqueDeviceCode,
			deviceCode,
			DuplicateDeviceCode,
		),
		eventstore.NewAddEventUniqueConstraint(
			UniqueUserCode,
			userCode,
			DuplicateUserCode,
		),
	}
}

func NewRemoveUniqueConstraints(deviceCode, userCode string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewRemoveUniqueConstraint(
			UniqueDeviceCode,
			deviceCode,
		),
		eventstore.NewRemoveUniqueConstraint(
			UniqueUserCode,
			userCode,
		),
	}
}
