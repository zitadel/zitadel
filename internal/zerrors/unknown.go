package zerrors

import "fmt"

func ThrowUnknown(parent error, id, message string) error {
	return newZitadelError(KindUnknown, parent, id, message)
}

func ThrowUnknownf(parent error, id, format string, a ...any) error {
	return newZitadelError(KindUnknown, parent, id, fmt.Sprintf(format, a...))
}

func IsUnknown(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnknown
}
