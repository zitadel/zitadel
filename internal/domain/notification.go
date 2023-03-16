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
	NotificationProviderStateActive
	NotificationProviderStateRemoved

	notificationProviderCount
)

func (s NotificationProviderState) Exists() bool {
	return s == NotificationProviderStateActive
}

type NotificationProviderType int32

const (
	NotificationProviderTypeFile NotificationProviderType = iota
	NotificationProviderTypeLog

	notificationProviderTypeCount
)
