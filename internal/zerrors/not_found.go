package zerrors

import "fmt"

func ThrowNotFound(parent error, id, message string) error {
	return CreateZitadelError(KindNotFound, parent, id, message, 1)
}

func ThrowNotFoundSlug(parent error, slug Slug, message string) error {
	return CreateZitadelErrorSlug(KindNotFound, parent, slug, message, 1)
}

func ThrowNotFoundSlugWithDetails(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelErrorSlug(KindNotFound, parent, slug, message, 1).WithDetails(details)
}

func ThrowNotFoundf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindNotFound, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsNotFound(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindNotFound
}
