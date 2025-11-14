package errors

type WrongPasswordError struct {
	FailedAttempts int32
}

func (wpe *WrongPasswordError) Error() string {
	return ""
}
