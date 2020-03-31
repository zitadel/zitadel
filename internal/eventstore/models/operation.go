package models

type Operation int32

const (
	Operation_Equals Operation = iota
	Operation_Greater
	Operation_Less
	Operation_In
)
