package zerrors

import "fmt"

func ThrowNotFound(parent error, id, message string) error {
	return newZitadelError(KindNotFound, parent, id, message)
}

func ThrowNotFoundf(parent error, id, format string, a ...any) error {
	return newZitadelError(KindNotFound, parent, id, fmt.Sprintf(format, a...))
}

func IsNotFound(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindNotFound
}
