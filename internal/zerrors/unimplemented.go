package zerrors

import "fmt"

func ThrowUnimplemented(parent error, id, message string) error {
	return CreateZitadelError(KindUnimplemented, parent, id, message)
}

func ThrowUnimplementedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindUnimplemented, parent, id, fmt.Sprintf(format, a...))
}

func IsUnimplemented(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnimplemented
}
