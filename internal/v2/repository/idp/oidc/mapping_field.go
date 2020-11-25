package oidc

type MappingField int32

const (
	MappingFieldPreferredLoginName MappingField = iota + 1
	MappingFieldEmail
	// count is for validation purposes
	mappingFieldCount
)

func (f MappingField) Valid() bool {
	return f > 0 && f < mappingFieldCount
}
