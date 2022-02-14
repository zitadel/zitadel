package domain

type NotificationType int32

const (
	NotificationTypeEmail NotificationType = iota
	NotificationTypeSms

	notificationCount
)

func (f NotificationType) Valid() bool {
	return f >= 0 && f < notificationCount
}

type NotificationProviderState int32

const (
	NotificationProviderStateUnspecified NotificationProviderState = iota
	NotificationProviderStateEnabled
	NotificationProviderStateDisabled
	NotificationProviderStateRemoved

	notificationProviderCount
)

func (s NotificationProviderState) Exists() bool {
	return s == NotificationProviderStateEnabled || s == NotificationProviderStateDisabled
}
