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
	RemoveMFAU2F(ctx context.Context, userID, webAuthNTokenID string) error

	GetPasswordless(ctx context.Context, id string) ([]*model.WebAuthNToken, error)
	AddPasswordless(ctx context.Context, id string) (*model.WebAuthNToken, error)
	VerifyPasswordlessSetup(ctx context.Context, userID, tokenName, userAgentID string, credentialData []byte) error
	RemovePasswordless(ctx context.Context, userID, webAuthNTokenID string) error

	ChangeUsername(ctx context.Context, userID, username string) error

	SignOut(ctx context.Context, agentID string) error

	UserByID(ctx context.Context, userID string) (*model.UserView, error)

	MachineKeyByID(ctx context.Context, keyID string) (*model.MachineKeyView, error)
}

type myUserRepo interface {
	MyUser(ctx context.Context) (*model.UserView, error)

	MyProfile(ctx context.Context) (*model.Profile, error)

	MyEmail(ctx context.Context) (*model.Email, error)
	ChangeMyEmail(ctx context.Context, email *model.Email) (*model.Email, error)
	VerifyMyEmail(ctx context.Context, code string) error
	ResendMyEmailVerificationMail(ctx context.Context) error

	MyPhone(ctx context.Context) (*model.Phone, error)
	ChangeMyPhone(ctx context.Context, phone *model.Phone) (*model.Phone, error)
	RemoveMyPhone(ctx context.Context) error
	VerifyMyPhone(ctx context.Context, code string) error
	ResendMyPhoneVerificationCode(ctx context.Context) error

	MyAddress(ctx context.Context) (*model.Address, error)
	ChangeMyAddress(ctx context.Context, address *model.Address) (*model.Address, error)

	ChangeMyPassword(ctx context.Context, old, new string) error

	SearchMyExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error)
	AddMyExternalIDP(ctx context.Context, externalIDP *model.ExternalIDP) (*model.ExternalIDP, error)
	RemoveMyExternalIDP(ctx context.Context, externalIDP *model.ExternalIDP) error

	MyUserMFAs(ctx context.Context) ([]*model.MultiFactor, error)
	AddMyMFAOTP(ctx context.Context) (*model.OTP, error)
	VerifyMyMFAOTPSetup(ctx context.Context, code string) error
	RemoveMyMFAOTP(ctx context.Context) error

	AddMyMFAU2F(ctx context.Context) (*model.WebAuthNToken, error)
	VerifyMyMFAU2FSetup(ctx context.Context, tokenName string, data []byte) error
	RemoveMyMFAU2F(ctx context.Context, webAuthNTokenID string) error

	GetMyPasswordless(ctx context.Context) ([]*model.WebAuthNToken, error)
	AddMyPasswordless(ctx context.Context) (*model.WebAuthNToken, error)
	VerifyMyPasswordlessSetup(ctx context.Context, tokenName string, data []byte) error
	RemoveMyPasswordless(ctx context.Context, webAuthNTokenID string) error

	MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error)
}
