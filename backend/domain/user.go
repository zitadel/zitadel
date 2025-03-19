package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/handler"
	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type User struct {
	db database.Pool

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

type VerifyEmail struct {
	UserID string
	Code   string
	Alg    crypto.EncryptionAlgorithm

	client          database.QueryExecutor
	config          *crypto.GeneratorConfig
	gen             crypto.Generator
	code            *repository.EmailVerificationCode
	verificationErr error
}

func (u *User) VerifyEmail(ctx context.Context, in *VerifyEmail) error {
	_, err := handler.Deferrable(
		func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, _ func(context.Context, error) error, err error) {
			client, err := u.db.Acquire(ctx)
			if err != nil {
				return nil, nil, err
			}
			in.client = client
			return in, func(ctx context.Context, _ error) error { return client.Release(ctx) }, err
		},
		handler.Chains(
			func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, err error) {
				in.config, err = u.secretGenerator.GeneratorConfigByType(ctx, in.client, domain.SecretGeneratorTypeVerifyEmailCode)
				return in, err
			},
			func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, err error) {
				in.gen = crypto.NewEncryptionGenerator(*in.config, in.Alg)
				return in, nil
			},
			handler.Deferrable(
				func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, _ func(context.Context, error) error, err error) {
					client := in.client
					tx, err := in.client.(database.Client).Begin(ctx, nil)
					if err != nil {
						return nil, nil, err
					}
					in.client = tx
					return in, func(ctx context.Context, err error) error {
						err = tx.End(ctx, err)
						if err != nil {
							return err
						}
						in.client = client
						return nil
					}, err
				},
				handler.Chains(
					func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, err error) {
						in.code, err = u.user.EmailVerificationCode(ctx, in.client, in.UserID)
						return in, err
					},
					func(ctx context.Context, in *VerifyEmail) (*VerifyEmail, error) {
						in.verificationErr = crypto.VerifyCode(in.code.CreatedAt, in.code.Expiry, in.code.Code, in.Code, in.gen.Alg())
						return in, nil
					},
					handler.HandleIf(
						func(in *VerifyEmail) bool {
							return in.verificationErr == nil
						},
						func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, err error) {
							return in, u.user.EmailVerificationSucceeded(ctx, in.client, in.code)
						},
					),
					handler.HandleIf(
						func(in *VerifyEmail) bool {
							return in.verificationErr != nil
						},
						func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, err error) {
							return in, u.user.EmailVerificationFailed(ctx, in.client, in.code)
						},
					),
				),
			),
		),
	)(ctx, in)
	return err
}
