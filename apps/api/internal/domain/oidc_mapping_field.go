package domain

type OIDCMappingField int32

const (
	OIDCMappingFieldUnspecified OIDCMappingField = iota
	OIDCMappingFieldPreferredLoginName
	OIDCMappingFieldEmail
	// count is for validation purposes
	oidcMappingFieldCount
)

func (f OIDCMappingField) Valid() bool {
	return f > 0 && f < oidcMappingFieldCount
}
