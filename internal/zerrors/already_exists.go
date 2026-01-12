package zerrors

import "fmt"

func ThrowAlreadyExists(parent error, id, message string) error {
	return CreateZitadelError(KindAlreadyExists, parent, id, message)
}

func ThrowAlreadyExistsf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindAlreadyExists, parent, id, fmt.Sprintf(format, a...))
}

func IsErrorAlreadyExists(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindAlreadyExists
}
