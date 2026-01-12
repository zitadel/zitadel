package zerrors

import "fmt"

func ThrowUnauthenticated(parent error, id, message string) error {
	return CreateZitadelError(KindUnauthenticated, parent, id, message)
}

func ThrowUnauthenticatedf(parent error, id, format string, a ...any) error {
	return CreateZitadelError(KindUnauthenticated, parent, id, fmt.Sprintf(format, a...))
}

func IsUnauthenticated(err error) bool {
	zitadelErr, ok := AsZitadelError(err)
	return ok && zitadelErr.Kind == KindUnauthenticated
}
