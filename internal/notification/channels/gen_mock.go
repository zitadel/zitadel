package channels

//go:generate mockgen -package mock -destination ./mock/provider.mock.go github.com/caos/zitadel/internal/notification/providers NotificationProvider
//go:generate mockgen -package mock -destination ./mock/message.mock.go github.com/caos/zitadel/internal/notification/providers Message
