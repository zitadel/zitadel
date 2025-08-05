package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
)

func UpdateMyEmailToDomain(ctx context.Context, email *auth.SetMyEmailRequest) *domain.Email {
	return &domain.Email{
		ObjectRoot:   ctxToObjectRoot(ctx),
		EmailAddress: domain.EmailAddress(email.Email),
	}
}
