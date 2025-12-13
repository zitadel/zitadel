package zerrors

import "fmt"

func ThrowInternal(parent error, id, message string) error {
	return CreateZitadelError(KindInternal, parent, id, message)
}

func ThrowInternalf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindInternal, parent, id, fmt.Sprintf(format, a...))
}

func IsInternal(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindInternal
}
