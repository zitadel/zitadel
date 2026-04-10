package zerrors

import "fmt"

func ThrowUnavailableError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindUnavailable, parent, string(slug), message, 1).WithDetails(details)
}

func ThrowUnavailable(parent error, id, message string) error {
	return CreateZitadelError(KindUnavailable, parent, id, message, 1)
}

func ThrowUnavailablef(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindUnavailable, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsUnavailable(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnavailable
}
