package eventstore

import (
	"context"
	"github.com/caos/logging"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"strings"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IamRepository struct {
	SearchLimit uint64
	*iam_es.IamEventstore
	View           *admin_view.View
	SystemDefaults systemdefaults.SystemDefaults
	Roles          []string
}

func (repo *IamRepository) IamMemberByID(ctx context.Context, orgID, userID string) (*iam_model.IamMemberView, error) {
	member, err := repo.View.IamMemberByIDs(orgID, userID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IamMemberToModel(member), nil
}

func (repo *IamRepository) AddIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error) {
	member.AggregateID = repo.SystemDefaults.IamID
	return repo.IamEventstore.AddIamMember(ctx, member)
}

func (repo *IamRepository) ChangeIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error) {
	member.AggregateID = repo.SystemDefaults.IamID
	return repo.IamEventstore.ChangeIamMember(ctx, member)
}

func (repo *IamRepository) RemoveIamMember(ctx context.Context, userID string) error {
	member := iam_model.NewIamMember(repo.SystemDefaults.IamID, userID)
	return repo.IamEventstore.RemoveIamMember(ctx, member)
}

func (repo *IamRepository) SearchIamMembers(ctx context.Context, request *iam_model.IamMemberSearchRequest) (*iam_model.IamMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIamMemberSequence()
	logging.Log("EVENT-Slkci").OnError(err).Warn("could not read latest iam sequence")
	members, count, err := repo.View.SearchIamMembers(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IamMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      iam_es_model.IamMembersToModel(members),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *IamRepository) IdpConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IdpConfigView, error) {
	idp, err := repo.View.IdpConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IdpConfigViewToModel(idp), nil
}
func (repo *IamRepository) AddOidcIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error) {
	idp.AggregateID = repo.SystemDefaults.IamID
	return repo.IamEventstore.AddIdpConfiguration(ctx, idp)
}

func (repo *IamRepository) ChangeIdpConfig(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error) {
	idp.AggregateID = repo.SystemDefaults.IamID
	return repo.IamEventstore.ChangeIdpConfiguration(ctx, idp)
}

func (repo *IamRepository) DeactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error) {
	return repo.IamEventstore.DeactivateIdpConfiguration(ctx, repo.SystemDefaults.IamID, idpConfigID)
}

func (repo *IamRepository) ReactivateIdpConfig(ctx context.Context, idpConfigID string) (*iam_model.IdpConfig, error) {
	return repo.IamEventstore.ReactivateIdpConfiguration(ctx, repo.SystemDefaults.IamID, idpConfigID)
}

func (repo *IamRepository) RemoveIdpConfig(ctx context.Context, idpConfigID string) error {
	//TODO: Remove from all policies and users
	idp := iam_model.NewIdpConfig(repo.SystemDefaults.IamID, idpConfigID)
	return repo.IamEventstore.RemoveIdpConfiguration(ctx, idp)
}

func (repo *IamRepository) ChangeOidcIdpConfig(ctx context.Context, oidcConfig *iam_model.OidcIdpConfig) (*iam_model.OidcIdpConfig, error) {
	oidcConfig.AggregateID = repo.SystemDefaults.IamID
	return repo.IamEventstore.ChangeIdpOidcConfiguration(ctx, oidcConfig)
}

func (repo *IamRepository) SearchIdpConfigs(ctx context.Context, request *iam_model.IdpConfigSearchRequest) (*iam_model.IdpConfigSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIdpConfigSequence()
	logging.Log("EVENT-Dk8si").OnError(err).Warn("could not read latest idp config sequence")
	idps, count, err := repo.View.SearchIdpConfigs(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IdpConfigSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      iam_es_model.IdpConfigViewsToModel(idps),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *IamRepository) GetIamMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "IAM") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}
