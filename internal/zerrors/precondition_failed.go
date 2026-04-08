package zerrors

import "fmt"

func ThrowPreconditionFailedError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindPreconditionFailed, parent, string(slug), message, 1).WithDetails(details)
}

// Deprecated: use ThrowPreconditionFailedError instead
func ThrowPreconditionFailed(parent error, id, message string) error {
	return CreateZitadelError(KindPreconditionFailed, parent, id, message, 1)
}

// Deprecated: use ThrowPreconditionFailedError instead
func ThrowPreconditionFailedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindPreconditionFailed, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsPreconditionFailed(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindPreconditionFailed
}
