package user

type NotificationType int32

const (
	NotificationTypeEmail NotificationType = iota
	NotificationTypeSms
)
