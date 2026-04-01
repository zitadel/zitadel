package zerrors

import "fmt"

func ThrowInternal(parent error, id, message string) error {
	return CreateZitadelError(KindInternal, parent, id, message, 1)
}

func ThrowInternalSlug(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelErrorSlug(KindInternal, parent, slug, message, 1).WithDetails(details)
}

func ThrowInternalf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindInternal, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsInternal(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindInternal
}
