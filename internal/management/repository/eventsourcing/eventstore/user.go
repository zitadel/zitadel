package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	policy_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_types "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserRepo struct {
	SearchLimit  uint64
	UserEvents   *usr_event.UserEventstore
	PolicyEvents *policy_event.PolicyEventstore
	View         *view.View
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (project *usr_model.User, err error) {
	return repo.UserEvents.UserByID(ctx, id)
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	policy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return repo.UserEvents.CreateUser(ctx, user, policy)
}

func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	policyResourceOwner := auth.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	policy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, policyResourceOwner)
	if err != nil {
		return nil, err
	}
	return repo.UserEvents.RegisterUser(ctx, user, policy, resourceOwner)
}

func (repo *UserRepo) DeactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.DeactivateUser(ctx, id)
}

func (repo *UserRepo) ReactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.ReactivateUser(ctx, id)
}

func (repo *UserRepo) LockUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.LockUser(ctx, id)
}

func (repo *UserRepo) UnlockUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.UnlockUser(ctx, id)
}

func (repo *UserRepo) SearchUsers(ctx context.Context, request *usr_model.UserSearchRequest) (*usr_model.UserSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	projects, count, err := repo.View.SearchUsers(request)
	if err != nil {
		return nil, err
	}
	return &usr_model.UserSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.UsersToModel(projects),
	}, nil
}

func (repo *UserRepo) UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64) (*usr_model.UserChanges, error) {
	changes, err := repo.UserEvents.UserChanges(ctx, usr_types.UserAggregate, id, lastSequence, limit)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func (repo *UserRepo) GetGlobalUserByEmail(ctx context.Context, email string) (*usr_model.UserView, error) {
	user, err := repo.View.GetGlobalUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return model.UserToModel(user), nil
}

func (repo *UserRepo) IsUserUnique(ctx context.Context, userName, email string) (bool, error) {
	return repo.View.IsUserUnique(userName, email)
}

func (repo *UserRepo) UserMfas(ctx context.Context, userID string) ([]*usr_model.MultiFactor, error) {
	return repo.View.UserMfas(userID)
}

func (repo *UserRepo) SetOneTimePassword(ctx context.Context, password *usr_model.Password) (*usr_model.Password, error) {
	policy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return repo.UserEvents.SetOneTimePassword(ctx, policy, password)
}

func (repo *UserRepo) RequestSetPassword(ctx context.Context, id string, notifyType usr_model.NotificationType) error {
	return repo.UserEvents.RequestSetPassword(ctx, id, notifyType)
}

func (repo *UserRepo) ProfileByID(ctx context.Context, userID string) (*usr_model.Profile, error) {
	return repo.UserEvents.ProfileByID(ctx, userID)
}

func (repo *UserRepo) ChangeProfile(ctx context.Context, profile *usr_model.Profile) (*usr_model.Profile, error) {
	return repo.UserEvents.ChangeProfile(ctx, profile)
}

func (repo *UserRepo) EmailByID(ctx context.Context, userID string) (*usr_model.Email, error) {
	return repo.UserEvents.EmailByID(ctx, userID)
}

func (repo *UserRepo) ChangeEmail(ctx context.Context, email *usr_model.Email) (*usr_model.Email, error) {
	return repo.UserEvents.ChangeEmail(ctx, email)
}

func (repo *UserRepo) CreateEmailVerificationCode(ctx context.Context, userID string) error {
	return repo.UserEvents.CreateEmailVerificationCode(ctx, userID)
}

func (repo *UserRepo) PhoneByID(ctx context.Context, userID string) (*usr_model.Phone, error) {
	return repo.UserEvents.PhoneByID(ctx, userID)
}

func (repo *UserRepo) ChangePhone(ctx context.Context, email *usr_model.Phone) (*usr_model.Phone, error) {
	return repo.UserEvents.ChangePhone(ctx, email)
}

func (repo *UserRepo) CreatePhoneVerificationCode(ctx context.Context, userID string) error {
	return repo.UserEvents.CreatePhoneVerificationCode(ctx, userID)
}

func (repo *UserRepo) AddressByID(ctx context.Context, userID string) (*usr_model.Address, error) {
	return repo.UserEvents.AddressByID(ctx, userID)
}

func (repo *UserRepo) ChangeAddress(ctx context.Context, address *usr_model.Address) (*usr_model.Address, error) {
	return repo.UserEvents.ChangeAddress(ctx, address)
}
