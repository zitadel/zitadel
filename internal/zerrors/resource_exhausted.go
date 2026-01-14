package zerrors

import "fmt"

func ThrowResourceExhausted(parent error, id, message string) error {
	return CreateZitadelError(KindResourceExhausted, parent, id, message)
}

func ThrowResourceExhaustedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindResourceExhausted, parent, id, fmt.Sprintf(format, a...))
}

func IsResourceExhausted(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindResourceExhausted
}
