package zerrors

import "fmt"

func ThrowDeadlineExceeded(parent error, id, message string) error {
	return CreateZitadelError(KindDeadlineExceeded, parent, id, message, 1)
}

func ThrowDeadlineExceededf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindDeadlineExceeded, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsDeadlineExceeded(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindDeadlineExceeded
}
