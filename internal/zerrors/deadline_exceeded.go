package zerrors

import "fmt"

func ThrowDeadlineExceededError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindDeadlineExceeded, parent, string(slug), message, 1).WithDetails(details)
}

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
