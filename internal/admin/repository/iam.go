package repository

import (
	"context"
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

	GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)
	AddDefaultMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error)
	ChangeDefaultMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error)

	GetDefaultMailTexts(ctx context.Context) (*iam_model.MailTextsView, error)
	GetDefaultMailText(ctx context.Context, textType string, language string) (*iam_model.MailTextView, error)
	AddDefaultMailText(ctx context.Context, mailText *iam_model.MailText) (*iam_model.MailText, error)
	ChangeDefaultMailText(ctx context.Context, policy *iam_model.MailText) (*iam_model.MailText, error)

	GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)

	GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error)

	GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error)

	GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)
}
