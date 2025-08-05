package authz

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
)

// UserIDInCTX checks if the userID
// equals the authenticated user in the context.
func UserIDInCTX(ctx context.Context, userID string) error {
	if GetCtxData(ctx).UserID != userID {
		return zerrors.ThrowPermissionDenied(nil, "AUTH-Bohd2", "Errors.User.UserIDWrong")
	}
	return nil
}
