package models

type Operation int32

const (
	Operation_Equals Operation = 1 + iota
	Operation_Greater
	Operation_Less
	Operation_In
)
