package repository

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	usr_model "github.com/caos/zitadel/internal/user/model"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	SearchIAMMembers(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error)

	GetIAMMemberRoles() []string

	SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)

	GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error)
	SearchDefaultIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error)
	SearchDefaultSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error)
	SearchDefaultMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error)

	IDPProvidersByIDPConfigID(ctx context.Context, idpConfigID string) ([]*iam_model.IDPProviderView, error)
	ExternalIDPsByIDPConfigID(ctx context.Context, idpConfigID string) ([]*usr_model.ExternalIDPView, error)
	ExternalIDPsByIDPConfigIDFromDefaultPolicy(ctx context.Context, idpConfigID string) ([]*usr_model.ExternalIDPView, error)

	GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetDefaultPreviewLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)

	GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)

	GetDefaultMessageTexts(ctx context.Context) (*iam_model.MessageTextsView, error)
	GetDefaultMessageText(ctx context.Context, textType string, language string) (*iam_model.MessageTextView, error)
	GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)

	GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)

	GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)

	GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)

	GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)
}
