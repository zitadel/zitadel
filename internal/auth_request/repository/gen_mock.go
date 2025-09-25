package repository

//go:generate go tool mockgen -package mock -destination ./mock/repository.mock.go github.com/zitadel/zitadel/internal/auth_request/repository AuthRequestCache
