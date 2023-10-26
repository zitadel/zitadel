package handlers

//go:generate mockgen -package mock -destination ./mock/queries.mock.go github.com/zitadel/zitadel/internal/notification/handlers Queries
//go:generate mockgen -package mock -destination ./mock/commands.mock.go github.com/zitadel/zitadel/internal/notification/handlers Commands
