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

func NewAddUniqueConstraints(clientID, deviceCode, userCode string) []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{
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

func NewRemoveUniqueConstraints(clientID, deviceCode, userCode string) []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{
		eventstore.NewRemoveEventUniqueConstraint(
			UniqueDeviceCode,
			deviceCodeUniqueField(clientID, deviceCode),
		),
		eventstore.NewRemoveEventUniqueConstraint(
			UniqueUserCode,
			userCode,
		),
	}
}
