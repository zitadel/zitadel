package feature

import "slices"

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
	KeyImprovedPerformance
	KeyWebKey
	KeyDebugOIDCParentError
	KeyOIDCSingleV1SessionTermination
	KeyDisableUserTokenEvent
	KeyEnableBackChannelLogout
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
	LoginDefaultOrg                 bool                      `json:"login_default_org,omitempty"`
	TriggerIntrospectionProjections bool                      `json:"trigger_introspection_projections,omitempty"`
	LegacyIntrospection             bool                      `json:"legacy_introspection,omitempty"`
	UserSchema                      bool                      `json:"user_schema,omitempty"`
	TokenExchange                   bool                      `json:"token_exchange,omitempty"`
	Actions                         bool                      `json:"actions,omitempty"`
	ImprovedPerformance             []ImprovedPerformanceType `json:"improved_performance,omitempty"`
	WebKey                          bool                      `json:"web_key,omitempty"`
	DebugOIDCParentError            bool                      `json:"debug_oidc_parent_error,omitempty"`
	OIDCSingleV1SessionTermination  bool                      `json:"oidc_single_v1_session_termination,omitempty"`
	DisableUserTokenEvent           bool                      `json:"disable_user_token_event,omitempty"`
	EnableBackChannelLogout         bool                      `json:"enable_back_channel_logout,omitempty"`
}

type ImprovedPerformanceType int32

const (
	ImprovedPerformanceTypeUnknown = iota
	ImprovedPerformanceTypeOrgByID
	ImprovedPerformanceTypeProjectGrant
	ImprovedPerformanceTypeProject
	ImprovedPerformanceTypeUserGrant
	ImprovedPerformanceTypeOrgDomainVerified
)

func (f Features) ShouldUseImprovedPerformance(typ ImprovedPerformanceType) bool {
	return slices.Contains(f.ImprovedPerformance, typ)
}
