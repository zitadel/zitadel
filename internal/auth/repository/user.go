package repository

import (
	"context"

	org_model "github.com/caos/zitadel/internal/org/model"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	Register(ctx context.Context, user *model.User, member *org_model.OrgMember, resourceOwner string) (*model.User, error)

	myUserRepo
	SkipMfaInit(ctx context.Context, userID string) error

	RequestPasswordReset(ctx context.Context, username string) error
	SetPassword(ctx context.Context, userID, code, password string) error
	ChangePassword(ctx context.Context, userID, old, new string) error

	VerifyEmail(ctx context.Context, userID, code string) error
	ResendEmailVerificationMail(ctx context.Context, userID string) error

	AddMfaOTP(ctx context.Context, userID string) (*model.OTP, error)
	VerifyMfaOTPSetup(ctx context.Context, userID, code string) error

	SignOut(ctx context.Context, agentID string) error

	UserByID(ctx context.Context, userID string) (*model.UserView, error)
}

type myUserRepo interface {
	MyUser(ctx context.Context) (*model.UserView, error)

	MyProfile(ctx context.Context) (*model.Profile, error)
	ChangeMyProfile(ctx context.Context, profile *model.Profile) (*model.Profile, error)

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

	MyUserMfas(ctx context.Context) ([]*model.MultiFactor, error)
	AddMyMfaOTP(ctx context.Context) (*model.OTP, error)
	VerifyMyMfaOTPSetup(ctx context.Context, code string) error
	RemoveMyMfaOTP(ctx context.Context) error

	MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error)
}
