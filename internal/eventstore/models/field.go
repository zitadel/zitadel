package models

type Field int32

const (
	Field_AggregateType Field = iota
	Field_AggregateID
	Field_LatestSequence
	Field_ResourceOwner
	Field_ModifierService
	Field_ModifierUser
	Field_ModifierTenant
)
