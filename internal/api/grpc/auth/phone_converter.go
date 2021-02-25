package auth

import (
	"context"

	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func UpdateMyPhoneToDomain(ctx context.Context, phone *auth.SetMyPhoneRequest) *domain.Phone {
	return &domain.Phone{
		ObjectRoot:  ctxToObjectRoot(ctx),
		PhoneNumber: phone.Phone,
	}
}
