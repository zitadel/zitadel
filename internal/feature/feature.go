package feature

//go:generate enumer -type Key -transform snake -trimprefix Key
type Key int

const (
	KeyUnspecified Key = iota
	KeyLoginDefaultOrg
	KeyTriggerIntrospectionProjections
	KeyLegacyIntrospection
	KeyUserSchema
	KeyTokenExchange
	KeyActions
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
	UserSchema                      bool `json:"user_schema,omitempty"`
	TokenExchange                   bool `json:"token_exchange,omitempty"`
	Actions                         bool `json:"actions,omitempty"`
}
