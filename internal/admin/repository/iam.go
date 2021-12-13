package repository

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	Languages(ctx context.Context) ([]language.Tag, error)

	GetIAMMemberRoles() []string

	GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)

	GetDefaultMessageText(ctx context.Context, textType, language string) (*domain.CustomMessageText, error)
	GetCustomMessageText(ctx context.Context, textType string, language string) (*domain.CustomMessageText, error)
	GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)
	GetCustomLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)
}
