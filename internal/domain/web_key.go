package domain

type WebKeyState int

const (
	WebKeyStateUndefined WebKeyState = iota
	WebKeyStateInactive
	WebKeyStateActive
	WebWeyStateRemoved
)
