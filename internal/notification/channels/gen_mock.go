package channels

//go:generate go tool mockgen -package mock -destination ./mock/channel.mock.go github.com/zitadel/zitadel/internal/notification/channels NotificationChannel
//go:generate go tool mockgen -package mock -destination ./mock/message.mock.go github.com/zitadel/zitadel/internal/notification/channels Message
