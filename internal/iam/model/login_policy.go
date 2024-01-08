package model

type PolicyState int32

const (
	PolicyStateActive PolicyState = iota
	PolicyStateRemoved
)
