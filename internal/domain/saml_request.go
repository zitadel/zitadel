package domain

type SAMLRequestState int

const (
	SAMLRequestStateUnspecified SAMLRequestState = iota
	SAMLRequestStateAdded
	SAMLRequestStateFailed
	SAMLRequestStateSucceeded
)
