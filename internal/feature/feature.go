package feature

//go:generate enumer -type Feature -transform snake
type Feature int

const (
	Unspecified Feature = iota
	LoginDefaultOrg
	TriggerIntrospectionProjections
	LegacyIntrospection
)

//go:generate enumer -type Level -transform snake -trimprefix Level
type Level int

const (
	LevelUnspecified Level = iota
	LevelSystem
	LevelInstance
	LevelOrg
	LevelProject
	LevelApp
	LevelUser
)

type Features struct {
	LoginDefaultOrg                 bool `json:"login_default_org,omitempty"`
	TriggerIntrospectionProjections bool `json:"trigger_introspection_projections,omitempty"`
	LegacyIntrospection             bool `json:"legacy_introspection,omitempty"`
}
