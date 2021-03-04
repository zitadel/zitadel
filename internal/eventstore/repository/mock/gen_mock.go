package mock

//go:generate mockgen -package mock -destination ./repository.mock.go github.com/caos/zitadel/internal/eventstore/repository Repository
