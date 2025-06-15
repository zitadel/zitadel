package feature

import (
	"net/url"
)

//go:generate enumer -type Key -transform snake -trimprefix Key
type Key int

const (
	KeyUnspecified Key = iota
	KeyLoginDefaultOrg
	KeyTriggerIntrospectionProjections
	KeyLegacyIntrospection
	KeyUserSchema
	KeyTokenExchange
	KeyActionsDeprecated
	// This key is kept for backwards compatibility, do not use in new code.
	KeyImprovedPerformance
	KeyWebKey
	KeyDebugOIDCParentError
	KeyOIDCSingleV1SessionTermination
	KeyDisableUserTokenEvent
	KeyEnableBackChannelLogout
	KeyLoginV2
	KeyPermissionCheckV2
	KeyConsoleUseV2UserApi
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
	LoginDefaultOrg                 bool    `json:"login_default_org,omitempty"`
	TriggerIntrospectionProjections bool    `json:"trigger_introspection_projections,omitempty"`
	LegacyIntrospection             bool    `json:"legacy_introspection,omitempty"`
	UserSchema                      bool    `json:"user_schema,omitempty"`
	TokenExchange                   bool    `json:"token_exchange,omitempty"`
	WebKey                          bool    `json:"web_key,omitempty"`
	DebugOIDCParentError            bool    `json:"debug_oidc_parent_error,omitempty"`
	OIDCSingleV1SessionTermination  bool    `json:"oidc_single_v1_session_termination,omitempty"`
	DisableUserTokenEvent           bool    `json:"disable_user_token_event,omitempty"`
	EnableBackChannelLogout         bool    `json:"enable_back_channel_logout,omitempty"`
	LoginV2                         LoginV2 `json:"login_v2,omitempty"`
	PermissionCheckV2               bool    `json:"permission_check_v2,omitempty"`
	ConsoleUseV2UserApi             bool    `json:"console_use_v2_user_api,omitempty"`
}

/* Note: do not generate the stringer or enumer for this type, is it breaks existing events */

type LoginV2 struct {
	Required bool     `json:"required,omitempty"`
	BaseURI  *url.URL `json:"base_uri,omitempty"`
}
