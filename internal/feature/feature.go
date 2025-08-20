package feature

import (
	"net/url"
	"slices"
)

//go:generate enumer -type Key -transform snake -trimprefix Key
type Key int

const (
	// Reserved: 2, 3, 6, 8, 11

	KeyUnspecified                    Key = 0
	KeyLoginDefaultOrg                Key = 1
	KeyUserSchema                     Key = 4
	KeyTokenExchange                  Key = 5
	KeyImprovedPerformance            Key = 7
	KeyDebugOIDCParentError           Key = 9
	KeyOIDCSingleV1SessionTermination Key = 10
	KeyEnableBackChannelLogout        Key = 12
	KeyLoginV2                        Key = 13
	KeyPermissionCheckV2              Key = 14
	KeyConsoleUseV2UserApi            Key = 15
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
	LoginDefaultOrg                bool                      `json:"login_default_org,omitempty"`
	UserSchema                     bool                      `json:"user_schema,omitempty"`
	TokenExchange                  bool                      `json:"token_exchange,omitempty"`
	ImprovedPerformance            []ImprovedPerformanceType `json:"improved_performance,omitempty"`
	DebugOIDCParentError           bool                      `json:"debug_oidc_parent_error,omitempty"`
	OIDCSingleV1SessionTermination bool                      `json:"oidc_single_v1_session_termination,omitempty"`
	DisableUserTokenEvent          bool                      `json:"disable_user_token_event,omitempty"`
	EnableBackChannelLogout        bool                      `json:"enable_back_channel_logout,omitempty"`
	LoginV2                        LoginV2                   `json:"login_v2,omitempty"`
	PermissionCheckV2              bool                      `json:"permission_check_v2,omitempty"`
	ConsoleUseV2UserApi            bool                      `json:"console_use_v2_user_api,omitempty"`
}

/* Note: do not generate the stringer or enumer for this type, is it breaks existing events */

type ImprovedPerformanceType int32

const (
	// Reserved: 1

	ImprovedPerformanceTypeUnspecified       ImprovedPerformanceType = 0
	ImprovedPerformanceTypeProjectGrant      ImprovedPerformanceType = 2
	ImprovedPerformanceTypeProject           ImprovedPerformanceType = 3
	ImprovedPerformanceTypeUserGrant         ImprovedPerformanceType = 4
	ImprovedPerformanceTypeOrgDomainVerified ImprovedPerformanceType = 5
)

func (f Features) ShouldUseImprovedPerformance(typ ImprovedPerformanceType) bool {
	return slices.Contains(f.ImprovedPerformance, typ)
}

type LoginV2 struct {
	Required bool     `json:"required,omitempty"`
	BaseURI  *url.URL `json:"base_uri,omitempty"`
}
