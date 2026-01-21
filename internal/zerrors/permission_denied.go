package zerrors

import "fmt"

func ThrowPermissionDenied(parent error, id, message string) error {
	return newZitadelError(KindPermissionDenied, parent, id, message)
}

func ThrowPermissionDeniedf(parent error, id, format string, a ...any) error {
	return newZitadelError(KindPermissionDenied, parent, id, fmt.Sprintf(format, a...))
}

func IsPermissionDenied(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindPermissionDenied
}
