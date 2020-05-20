package providers

type NotificationProvider interface {
	CanHandleMessage() bool
	HandleMessage() error
}
