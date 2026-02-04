package zerrors

import "fmt"

func ThrowInternal(parent error, id, message string) error {
	return CreateZitadelError(KindInternal, parent, id, message, 1)
}

func ThrowInternalf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindInternal, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsInternal(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindInternal
}
