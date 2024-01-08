package domain

type NotificationType int32

const (
	NotificationTypeEmail NotificationType = iota
	NotificationTypeSms

	notificationCount
)

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
