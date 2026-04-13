package domain

type SessionCheckError struct {
	Error          error
	FailedAttempts uint16
}
