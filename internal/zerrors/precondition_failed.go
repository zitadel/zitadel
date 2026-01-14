package zerrors

import "fmt"

func ThrowPreconditionFailed(parent error, id, message string) error {
	return CreateZitadelError(KindPreconditionFailed, parent, id, message)
}

func ThrowPreconditionFailedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindPreconditionFailed, parent, id, fmt.Sprintf(format, a...))
}

func IsPreconditionFailed(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindPreconditionFailed
}
