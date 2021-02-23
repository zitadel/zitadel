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
