package domain

type OrganizationSettingsState int32

const (
	OrganizationSettingsStateUnspecified OrganizationSettingsState = iota
	OrganizationSettingsStateActive
	OrganizationSettingsStateRemoved

	organizationSettingsStateCount
)

func (c OrganizationSettingsState) Valid() bool {
	return c >= 0 && c < organizationSettingsStateCount
}

func (s OrganizationSettingsState) Exists() bool {
	return s.Valid() && s != OrganizationSettingsStateUnspecified && s != OrganizationSettingsStateRemoved
}
