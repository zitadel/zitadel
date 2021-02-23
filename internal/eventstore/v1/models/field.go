package models

type Field int32

const (
	Field_AggregateType Field = 1 + iota
	Field_AggregateID
	Field_LatestSequence
	Field_ResourceOwner
	Field_EditorService
	Field_EditorUser
	Field_EventType
)
