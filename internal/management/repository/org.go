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

	GetOrgMemberRoles() []string

	GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)
	GetMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)

	GetDefaultMessageText(ctx context.Context, textType string, language string) (*domain.CustomMessageText, error)
	GetMessageText(ctx context.Context, orgID, textType, lang string) (*domain.CustomMessageText, error)

	GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)
	GetLoginTexts(ctx context.Context, orgID, lang string) (*domain.CustomLoginText, error)
}
