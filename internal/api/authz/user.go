package authz

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
)

// UserIdRequest is implemented by all v2 API requests
// that carry the user_id field.
type UserIdRequest interface {
	GetUserId() string
}

// RequestEqualsCTXUser checks if the v2 API request UserID
// equals the authenticated user in the context.
func RequestEqualsCTXUser(ctx context.Context, req UserIdRequest) error {
	if GetCtxData(ctx).UserID != req.GetUserId() {
		return errors.ThrowUnauthenticated(nil, "AUTH-Bohd2", "request user not equal to authenticated user")
	}
	return nil
}
