package zerrors

import (
	"fmt"
)

func ThrowPermissionDeniedError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindPermissionDenied, parent, string(slug), message, 1).WithDetails(details)
}

func ThrowPermissionDenied(parent error, id, message string) error {
	return CreateZitadelError(KindPermissionDenied, parent, id, message, 1)
}

func ThrowPermissionDeniedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindPermissionDenied, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsPermissionDenied(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindPermissionDenied
}
