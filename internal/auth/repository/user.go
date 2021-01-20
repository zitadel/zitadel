package repository

import (
	"context"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	myUserRepo

	RequestPasswordReset(ctx context.Context, username string) error

	ResendEmailVerificationMail(ctx context.Context, userID string) error

	VerifyInitCode(ctx context.Context, userID, code, password string) error

	AddMFAU2F(ctx context.Context, id string) (*model.WebAuthNToken, error)
	VerifyMFAU2FSetup(ctx context.Context, userID, tokenName, userAgentID string, credentialData []byte) error

	GetPasswordless(ctx context.Context, id string) ([]*model.WebAuthNToken, error)
	AddPasswordless(ctx context.Context, id string) (*model.WebAuthNToken, error)
	VerifyPasswordlessSetup(ctx context.Context, userID, tokenName, userAgentID string, credentialData []byte) error

	SignOut(ctx context.Context, agentID string) error

	UserByID(ctx context.Context, userID string) (*model.UserView, error)

	MachineKeyByID(ctx context.Context, keyID string) (*model.MachineKeyView, error)
}

type myUserRepo interface {
	MyUser(ctx context.Context) (*model.UserView, error)

	MyProfile(ctx context.Context) (*model.Profile, error)

	MyEmail(ctx context.Context) (*model.Email, error)

	MyPhone(ctx context.Context) (*model.Phone, error)

	MyAddress(ctx context.Context) (*model.Address, error)

	SearchMyExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error)

	MyUserMFAs(ctx context.Context) ([]*model.MultiFactor, error)

	GetMyPasswordless(ctx context.Context) ([]*model.WebAuthNToken, error)

	MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error)
}
