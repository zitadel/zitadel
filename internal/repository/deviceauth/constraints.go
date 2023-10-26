package deviceauth

import (
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	UniqueUserCode      = "user_code"
	UniqueDeviceCode    = "device_code"
	DuplicateUserCode   = "Errors.DeviceUserCode.AlreadyExists"
	DuplicateDeviceCode = "Errors.DeviceCode.AlreadyExists"
)

func deviceCodeUniqueField(clientID, deviceCode string) string {
	return strings.Join([]string{clientID, deviceCode}, ":")
}

func NewAddUniqueConstraints(clientID, deviceCode, userCode string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(
			UniqueDeviceCode,
			deviceCodeUniqueField(clientID, deviceCode),
			DuplicateDeviceCode,
		),
		eventstore.NewAddEventUniqueConstraint(
			UniqueUserCode,
			userCode,
			DuplicateUserCode,
		),
	}
}

func NewRemoveUniqueConstraints(clientID, deviceCode, userCode string) []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewRemoveUniqueConstraint(
			UniqueDeviceCode,
			deviceCodeUniqueField(clientID, deviceCode),
		),
		eventstore.NewRemoveUniqueConstraint(
			UniqueUserCode,
			userCode,
		),
	}
}
