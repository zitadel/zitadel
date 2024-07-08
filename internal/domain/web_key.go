package domain

type WebKeyState int

const (
	WebKeyStateUnspecified WebKeyState = iota
	WebKeyStateInactive
	WebKeyStateActive
	WebKeyStateRemoved
)
