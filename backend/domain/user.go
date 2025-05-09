package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
)

type User struct {
	defaults

	userCodeAlg     crypto.EncryptionAlgorithm
	user            userRepository
	secretGenerator secretGeneratorRepository
}

type UserRepositoryConstructor interface {
	NewUserExecutor(database.Executor) userRepository
	NewUserQuerier(database.Querier) userRepository
}

type userRepository interface {
	Create(ctx context.Context, tx database.Executor, user *repository.User) (*repository.User, error)
	ByID(ctx context.Context, querier database.Querier, id string) (*repository.User, error)

	EmailVerificationCode(ctx context.Context, client database.Querier, userID string) (*repository.EmailVerificationCode, error)
	EmailVerificationFailed(ctx context.Context, client database.Executor, code *repository.EmailVerificationCode) error
	EmailVerificationSucceeded(ctx context.Context, client database.Executor, code *repository.EmailVerificationCode) error
}

type secretGeneratorRepository interface {
	GeneratorConfigByType(ctx context.Context, client database.Querier, typ repository.SecretGeneratorType) (*crypto.GeneratorConfig, error)
}

func NewUser(db database.Pool) *User {
	b := &User{
		db:              db,
		user:            repository.NewUser(),
		secretGenerator: repository.NewSecretGenerator(),
	}

	return b
}
