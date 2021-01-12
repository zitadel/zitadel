package repository

import (
	"context"

	org_model "github.com/caos/zitadel/internal/org/model"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	Register(ctx context.Context, user *model.User, member *org_model.OrgMember, resourceOwner string) (*model.User, error)
	RegisterExternalUser(ctx context.Context, user *model.User, externalIDP *model.ExternalIDP, member *org_model.OrgMember, resourceOwner string) (*model.User, error)

	myUserRepo
	SkipMFAInit(ctx context.Context, userID string) error

	RequestPasswordReset(ctx context.Context, username string) error
	SetPassword(ctx context.Context, userID, code, password, userAgentID string) error
	ChangePassword(ctx context.Context, userID, old, new, userAgentID string) error

	VerifyEmail(ctx context.Context, userID, code string) error
	ResendEmailVerificationMail(ctx context.Context, userID string) error

	VerifyInitCode(ctx context.Context, userID, code, password string) error
	ResendInitVerificationMail(ctx context.Context, userID string) error

	AddMFAOTP(ctx context.Context, userID string) (*model.OTP, error)
	VerifyMFAOTPSetup(ctx context.Context, userID, code, userAgentID string) error

	AddMFAU2F(ctx context.Context, id string) (*model.WebAuthNToken, error)
	VerifyMFAU2FSetup(ctx context.Context, userID, tokenName, userAgentID string, credentialData []byte) error

	GetPasswordless(ctx context.Context, id string) ([]*model.WebAuthNToken, error)
	AddPasswordless(ctx context.Context, id string) (*model.WebAuthNToken, error)
	VerifyPasswordlessSetup(ctx context.Context, userID, tokenName, userAgentID string, credentialData []byte) error

	ChangeUsername(ctx context.Context, userID, username string) error

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

	AddMyMFAU2F(ctx context.Context) (*model.WebAuthNToken, error)
	VerifyMyMFAU2FSetup(ctx context.Context, tokenName string, data []byte) error

	GetMyPasswordless(ctx context.Context) ([]*model.WebAuthNToken, error)
	AddMyPasswordless(ctx context.Context) (*model.WebAuthNToken, error)
	VerifyMyPasswordlessSetup(ctx context.Context, tokenName string, data []byte) error

	MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error)
}
