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
	OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*org_model.OrgChanges, error)

	SearchMyOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error)

	GetOrgMemberRoles() []string

	SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)
	IDPConfigByID(ctx context.Context, id string) (*iam_model.IDPConfigView, error)

	SearchIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error)
	GetIDPProvidersByIDPConfigID(ctx context.Context, aggregateID, idpConfigID string) ([]*iam_model.IDPProviderView, error)

	GetDefaultMessageText(ctx context.Context, textType string, language string) (*domain.CustomMessageText, error)
	GetMessageText(ctx context.Context, orgID, textType, lang string) (*domain.CustomMessageText, error)

	GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)
	GetLoginTexts(ctx context.Context, orgID, lang string) (*domain.CustomLoginText, error)

	GetLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetPreviewLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetPreviewDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
}
