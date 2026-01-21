package zerrors

import "fmt"

func ThrowDeadlineExceeded(parent error, id, message string) error {
	return newZitadelError(KindDeadlineExceeded, parent, id, message)
}

func ThrowDeadlineExceededf(parent error, id, format string, a ...any) error {
	return newZitadelError(KindDeadlineExceeded, parent, id, fmt.Sprintf(format, a...))
}

func IsDeadlineExceeded(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindDeadlineExceeded
}
