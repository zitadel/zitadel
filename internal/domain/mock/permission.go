package permissionmock

import (
	"golang.org/x/net/context"

	"github.com/zitadel/zitadel/internal/domain"
)

// MockPermissionCheckErr returns a permission check function that will fail
// and return the input error
func MockPermissionCheckErr(err error) domain.PermissionCheck {
	return func(_ context.Context, _, _, _ string) error {
		return err
	}
}

// MockPermissionCheckOK returns a permission check function that will succeed
func MockPermissionCheckOK() domain.PermissionCheck {
	return func(_ context.Context, _, _, _ string) (err error) {
		return nil
	}
}
