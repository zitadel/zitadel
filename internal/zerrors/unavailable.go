package zerrors

import "fmt"

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
