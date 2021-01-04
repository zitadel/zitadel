package domain

type IDPConfigType int32

const (
	IDPConfigTypeOIDC IDPConfigType = iota
	IDPConfigTypeSAML

	//count is for validation
	idpConfigTypeCount
)

func (f IDPConfigType) Valid() bool {
	return f >= 0 && f < idpConfigTypeCount
}

type IDPConfigState int32

const (
	IDPConfigStateUnspecified IDPConfigState = iota
	IDPConfigStateActive
	IDPConfigStateInactive
	IDPConfigStateRemoved

	idpConfigStateCount
)

func (f IDPConfigState) Valid() bool {
	return f >= 0 && f < idpConfigStateCount
}

type IDPConfigStylingType int32

const (
	IDPConfigStylingTypeGoogle IDPConfigStylingType = iota + 1

	idpConfigStylingTypeCount
)

func (f IDPConfigStylingType) Valid() bool {
	return f >= 0 && f < idpConfigStylingTypeCount
}
