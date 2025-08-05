package execution

//go:generate mockgen -package mock -destination ./mock/queries.mock.go github.com/zitadel/zitadel/internal/execution Queries
//go:generate mockgen -package mock -destination ./mock/queue.mock.go github.com/zitadel/zitadel/internal/execution Queue
