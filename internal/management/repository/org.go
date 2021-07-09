package repository

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"

	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	Languages(ctx context.Context) ([]language.Tag, error)
	OrgByID(ctx context.Context, id string) (*org_model.OrgView, error)
	OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.OrgView, error)
	OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*org_model.OrgChanges, error)

	SearchMyOrgDomains(ctx context.Context, request *org_model.OrgDomainSearchRequest) (*org_model.OrgDomainSearchResponse, error)

	SearchMyOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error)

	GetOrgMemberRoles() []string

	SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)
	IDPConfigByID(ctx context.Context, id string) (*iam_model.IDPConfigView, error)

	GetMyOrgIamPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)

	GetLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	SearchIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error)
	GetIDPProvidersByIDPConfigID(ctx context.Context, aggregateID, idpConfigID string) ([]*iam_model.IDPProviderView, error)
	SearchSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error)
	SearchMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error)

	GetPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
	GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)

	GetPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)
	GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)

	GetPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)
	GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)

	GetPrivacyPolicy(ctx context.Context) (*iam_model.PrivacyPolicyView, error)
	GetDefaultPrivacyPolicy(ctx context.Context) (*iam_model.PrivacyPolicyView, error)

	GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)
	GetMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)

	GetDefaultMessageText(ctx context.Context, textType string, language string) (*domain.CustomMessageText, error)
	GetMessageText(ctx context.Context, orgID, textType, lang string) (*domain.CustomMessageText, error)

	GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)
	GetLoginTexts(ctx context.Context, orgID, lang string) (*domain.CustomLoginText, error)

	GetLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetPreviewLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetPreviewDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
}
