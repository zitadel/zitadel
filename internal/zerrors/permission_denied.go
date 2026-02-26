package zerrors

import "fmt"

func ThrowPermissionDenied(parent error, id, message string) error {
	return CreateZitadelError(KindPermissionDenied, parent, id, message, 1)
}

func ThrowPermissionDeniedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindPermissionDenied, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsPermissionDenied(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindPermissionDenied
}
