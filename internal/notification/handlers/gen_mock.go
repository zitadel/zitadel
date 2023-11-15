package handlers

//go:generate mockgen -package mock -destination ./mock/queries.mock.go github.com/zitadel/zitadel/v2/internal/notification/handlers Queries
//go:generate mockgen -package mock -destination ./mock/commands.mock.go github.com/zitadel/zitadel/v2/internal/notification/handlers Commands
