package zerrors

import "fmt"

func ThrowUnauthenticatedError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindUnauthenticated, parent, string(slug), message, 1).WithDetails(details)
}

// Deprecated: use ThrowUnauthenticatedError instead
func ThrowUnauthenticated(parent error, id, message string) error {
	return CreateZitadelError(KindUnauthenticated, parent, id, message, 1)
}

// Deprecated: use ThrowUnauthenticatedError instead
func ThrowUnauthenticatedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindUnauthenticated, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsUnauthenticated(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnauthenticated
}
