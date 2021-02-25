package auth

import (
	"context"

	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func UpdateAddressToDomain(ctx context.Context, address *auth.UpdateMyAddressRequest) *domain.Address {
	return &domain.Address{
		ObjectRoot:    ctxToObjectRoot(ctx),
		Country:       address.Country,
		Locality:      address.Locality,
		PostalCode:    address.PostalCode,
		Region:        address.Region,
		StreetAddress: address.StreetAddress,
	}
}
