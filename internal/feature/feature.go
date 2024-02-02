package feature

type Features struct {
	LoginDefaultOrg                 bool `json:"login_default_org,omitempty"`
	TriggerIntrospectionProjections bool `json:"trigger_introspection_projections,omitempty"`
	LegacyIntrospection             bool `json:"legacy_introspection,omitempty"`
}
