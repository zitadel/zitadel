package zerrors

import "fmt"

func ThrowUnavailable(parent error, id, message string) error {
	return newZitadelError(KindUnavailable, parent, id, message)
}

func ThrowUnavailablef(parent error, id, format string, a ...any) error {
	return newZitadelError(KindUnavailable, parent, id, fmt.Sprintf(format, a...))
}

func IsUnavailable(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnavailable
}
