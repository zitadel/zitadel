package zerrors

import "fmt"

func ThrowResourceExhaustedError(parent error, slug Slug, message string, details ErrorDetails) error {
	return CreateZitadelError(KindResourceExhausted, parent, string(slug), message, 1).WithDetails(details)
}

// Deprecated: use ThrowResourceExhaustedError instead
func ThrowResourceExhausted(parent error, id, message string) error {
	return CreateZitadelError(KindResourceExhausted, parent, id, message, 1)
}

// Deprecated: use ThrowResourceExhaustedError instead
func ThrowResourceExhaustedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindResourceExhausted, parent, id, fmt.Sprintf(format, a...), 1)
}

func IsResourceExhausted(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindResourceExhausted
}
