package errors

type WrongPasswordError struct {
	FailedAttempts int32
}

func (wpe *WrongPasswordError) Error() string {
	return ""
}

type LockDurationNotExceededError struct {
	RemainingTime int32
}

func (ldne *LockDurationNotExceededError) Error() string {
	return ""
}
