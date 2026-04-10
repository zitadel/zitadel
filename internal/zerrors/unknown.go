package zerrors

import "fmt"

func ThrowUnknownError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindUnknown, parent, string(slug), message, 1).WithDetails(details)
}

func ThrowUnknown(parent error, id, message string) error {
	return CreateZitadelError(KindUnknown, parent, id, message, 1)
}

func ThrowUnknownf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindUnknown, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsUnknown(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnknown
}
