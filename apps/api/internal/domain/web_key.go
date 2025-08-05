package domain

type WebKeyState int

const (
	WebKeyStateUnspecified WebKeyState = iota
	WebKeyStateInitial
	WebKeyStateActive
	WebKeyStateInactive
	WebKeyStateRemoved
)
