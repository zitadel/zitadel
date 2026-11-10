package zerrors

import "fmt"

func ThrowInvalidArgument(parent error, id, message string) error {
	return CreateZitadelError(KindInvalidArgument, parent, id, message, 1)
}

func ThrowInvalidArgumentf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindInvalidArgument, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsErrorInvalidArgument(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindInvalidArgument
}
