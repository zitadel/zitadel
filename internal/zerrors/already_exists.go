package zerrors

import "fmt"

func ThrowAlreadyExistsError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindAlreadyExists, parent, string(slug), message, 1).WithDetails(details)
}

func ThrowAlreadyExists(parent error, id, message string) error {
	return CreateZitadelError(KindAlreadyExists, parent, id, message, 1)
}

func ThrowAlreadyExistsf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindAlreadyExists, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsErrorAlreadyExists(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindAlreadyExists
}
