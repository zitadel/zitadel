package zerrors

import "fmt"

func ThrowUnimplementedError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindUnimplemented, parent, string(slug), message, 1).WithDetails(details)
}

func ThrowUnimplemented(parent error, id, message string) error {
	return CreateZitadelError(KindUnimplemented, parent, id, message, 1)
}

func ThrowUnimplementedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindUnimplemented, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsUnimplemented(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnimplemented
}
