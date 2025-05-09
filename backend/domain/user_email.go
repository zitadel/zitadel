package domain

import (
	"context"
	"text/template"

	"github.com/zitadel/zitadel/backend/handler"
	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type VerifyEmail struct {
	UserID string
	Code   string

	client          database.QueryExecutor
	config          *crypto.GeneratorConfig
	gen             crypto.Generator
	code            *repository.EmailVerificationCode
	verificationErr error
}

type SetEmail struct {
	*poolHandler[*SetEmail]

	UserID       string
	Email        string
	Verification handler.Handle[*SetEmail, *SetEmail]

	// config *crypto.GeneratorConfig
	gen crypto.Generator

	code      *crypto.CryptoValue
	plainCode string

	currentEmail string
}

func (u *User) WithEmailConfirmationURL(url template.Template) handler.Handle[*SetEmail, *SetEmail] {
	return handler.Chain(
		u.WithEmailReturnCode(),
		func(ctx context.Context, in *SetEmail) (out *SetEmail, err error) {
			// TODO: queue notification
			return in, nil
		},
	)
}

func (u *User) WithEmailReturnCode() handler.Handle[*SetEmail, *SetEmail] {
	return handler.Chains(
		handler.ErrFuncToHandle(
			func(ctx context.Context, in *SetEmail) (err error) {
				in.code, in.plainCode, err = crypto.NewCode(in.gen)
				return err
			},
		),
		handler.ErrFuncToHandle(
			func(ctx context.Context, in *SetEmail) (err error) {
				return u.user.SetEmailVerificationCode(ctx, in.poolHandler.client, in.UserID, in.code)
			},
		),
	)
}

func (u *User) WithEmailVerified() handler.Handle[*SetEmail, *SetEmail] {
	return handler.Chain(
		handler.ErrFuncToHandle(
			func(ctx context.Context, in *SetEmail) (err error) {
				return repository.SetEmailVerificationCode(ctx, in.poolHandler.client, in.UserID, in.code)
			},
		),
		handler.ErrFuncToHandle(
			func(ctx context.Context, in *SetEmail) (err error) {
				return u.user.EmailVerificationSucceeded(ctx, in.poolHandler.client, &repository.EmailVerificationCode{
					Code: in.code,
				})
			},
		),
	)
}

func (u *User) WithDefaultEmailVerification() handler.Handle[*SetEmail, *SetEmail] {
	return handler.Chain(
		u.WithEmailReturnCode(),
		func(ctx context.Context, in *SetEmail) (out *SetEmail, err error) {
			// TODO: queue notification
			return in, nil
		},
	)
}

func (u *User) SetEmailDifferent(ctx context.Context, in *SetEmail) (err error) {
	if in.Verification == nil {
		in.Verification = u.WithDefaultEmailVerification()
	}

	client, err := u.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer client.Release(ctx)

	config, err := u.secretGenerator.GeneratorConfigByType(ctx, client, domain.SecretGeneratorTypeVerifyEmailCode)
	if err != nil {
		return err
	}
	in.gen = crypto.NewEncryptionGenerator(*config, u.userCodeAlg)

	tx, err := client.Begin(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.End(ctx, err)

	user, err := u.user.ByID(ctx, tx, in.UserID)
	if err != nil {
		return err
	}

	if user.Email == in.Email {
		return nil
	}

	_, err = in.Verification(ctx, in)
	return err
}

func (u *User) SetEmail(ctx context.Context, in *SetEmail) error {
	_, err := handler.Chain(
		handler.HandleIf(
			func(in *SetEmail) bool {
				return in.Verification == nil
			},
			func(ctx context.Context, in *SetEmail) (*SetEmail, error) {
				in.Verification = u.WithDefaultEmailVerification()
				return in, nil
			},
		),
		handler.Deferrable(
			in.poolHandler.acquire,
			handler.Chains(
				func(ctx context.Context, in *SetEmail) (_ *SetEmail, err error) {
					config, err := u.secretGenerator.GeneratorConfigByType(ctx, in.poolHandler.client, domain.SecretGeneratorTypeVerifyEmailCode)
					if err != nil {
						return nil, err
					}
					in.gen = crypto.NewEncryptionGenerator(*config, u.userCodeAlg)
					return in, nil
				},
				handler.Deferrable(
					in.poolHandler.begin,
					handler.Chains(
						func(ctx context.Context, in *SetEmail) (*SetEmail, error) {
							// TODO: repository.EmailByUserID
							user, err := u.user.ByID(ctx, in.poolHandler.client, in.UserID)
							if err != nil {
								return nil, err
							}
							in.currentEmail = user.Email
							return in, nil
						},
						handler.SkipIf(
							func(in *SetEmail) bool {
								return in.currentEmail == in.Email
							},
							handler.Chains(
								func(ctx context.Context, in *SetEmail) (*SetEmail, error) {
									// TODO: repository.SetEmail
									return in, nil
								},
								in.Verification,
							),
						),
					),
				),
			),
		),
	)(ctx, in)
	return err
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
				in.gen = crypto.NewEncryptionGenerator(*in.config, u.userCodeAlg)
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
