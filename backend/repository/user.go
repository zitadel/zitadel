package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/handler"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/crypto"
)

type User struct {
	ID       string
	Username string
}

type UserOptions struct {
	cache *UserCache
}

type user struct {
	options[UserOptions]
	*UserOptions
}

func NewUser(opts ...Option[UserOptions]) *user {
	i := new(user)
	i.UserOptions = &i.options.custom

	for _, opt := range opts {
		opt(&i.options)
	}
	return i
}

func WithUserCache(c *UserCache) Option[UserOptions] {
	return func(i *options[UserOptions]) {
		i.custom.cache = c
	}
}

func (u *user) Create(ctx context.Context, tx database.Transaction, user *User) (*User, error) {
	return tracing.Wrap(u.tracer, "user.Create",
		handler.Chain(
			handler.Decorate(
				execute(tx).CreateUser,
				tracing.Decorate[*User, *User](u.tracer, tracing.WithSpanName("user.sql.Create")),
			),
			handler.Decorate(
				events(tx).CreateUser,
				tracing.Decorate[*User, *User](u.tracer, tracing.WithSpanName("user.event.Create")),
			),
		),
	)(ctx, user)
}

func (u *user) ByID(ctx context.Context, client database.Querier, id string) (*User, error) {
	return handler.SkipNext(
		handler.SkipNilHandler(u.cache,
			handler.ResFuncToHandle(u.cache.ByID),
		),
		handler.Chain(
			handler.Decorate(
				query(client).UserByID,
				tracing.Decorate[string, *User](u.tracer, tracing.WithSpanName("user.sql.ByID")),
			),
			handler.SkipNilHandler(u.custom.cache, handler.NoReturnToHandle(u.cache.Set)),
		),
	)(ctx, id)
}

type ChangeEmail struct {
	UserID string
	Email  string
	Opt    *ChangeEmailOption
}

type ChangeEmailOption struct {
	returnCode bool
	isVerified bool
	sendCode   bool
}

type ChangeEmailVerifiedOption struct {
	isVerified bool
}

type ChangeEmailReturnCodeOption struct {
	alg crypto.EncryptionAlgorithm
}

type ChangeEmailSendCodeOption struct {
	alg         crypto.EncryptionAlgorithm
	urlTemplate string
}

func (u *user) ChangeEmail(ctx context.Context, client database.Executor, change *ChangeEmail)
