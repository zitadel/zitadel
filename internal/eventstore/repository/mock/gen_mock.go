package mock

//go:generate go tool mockgen -package mock -destination ./repository.mock.go github.com/zitadel/zitadel/internal/eventstore Querier,Pusher
