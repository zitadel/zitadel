package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func newMockPermissionCheckAllowed() domain.PermissionCheck {
	return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return nil
	}
}

func newMockPermissionCheckNotAllowed() domain.PermissionCheck {
	return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied")
	}
}
