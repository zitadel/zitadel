package domain

import "github.com/zitadel/logging"

type IDPState int32

const (
	IDPStateUnspecified IDPState = iota
	IDPStateActive
	IDPStateInactive
	IDPStateRemoved
	IDPStateMigrated

	idpStateCount
)

func (s IDPState) Valid() bool {
	return s >= 0 && s < idpStateCount
}

func (s IDPState) Exists() bool {
	return s != IDPStateUnspecified && s != IDPStateRemoved && s != IDPStateMigrated
}

type IDPType int32

const (
	IDPTypeUnspecified IDPType = iota
	IDPTypeOIDC
	IDPTypeJWT
	IDPTypeOAuth
	IDPTypeLDAP
	IDPTypeAzureAD
	IDPTypeGitHub
	IDPTypeGitHubEnterprise
	IDPTypeGitLab
	IDPTypeGitLabSelfHosted
	IDPTypeGoogle
)

func (t IDPType) GetCSSClass() string {
	switch t {
	case IDPTypeGoogle:
		return "google"
	case IDPTypeGitHub,
		IDPTypeGitHubEnterprise:
		return "github"
	case IDPTypeGitLab,
		IDPTypeGitLabSelfHosted:
		return "gitlab"
	case IDPTypeAzureAD:
		return "azure"
	case IDPTypeUnspecified,
		IDPTypeOIDC,
		IDPTypeJWT,
		IDPTypeOAuth,
		IDPTypeLDAP:
		fallthrough
	default:
		return ""
	}
}

func IDPName(name string, idpType IDPType) string {
	if name != "" {
		return name
	}
	return idpType.DisplayName()
}

// DisplayName returns the name or a default
// to be used when always a name must be displayed (e.g. login)
func (t IDPType) DisplayName() string {
	switch t {
	case IDPTypeGitHub:
		return "GitHub"
	case IDPTypeGitLab:
		return "GitLab"
	case IDPTypeGoogle:
		return "Google"
	case IDPTypeUnspecified,
		IDPTypeOIDC,
		IDPTypeJWT,
		IDPTypeOAuth,
		IDPTypeLDAP,
		IDPTypeAzureAD,
		IDPTypeGitHubEnterprise,
		IDPTypeGitLabSelfHosted:
		fallthrough
	default:
		// we should never get here, so log it
		logging.Errorf("name of provider (type %d) is empty", t)
		return ""
	}
}

type IDPIntentState int32

const (
	IDPIntentStateUnspecified IDPIntentState = iota
	IDPIntentStateStarted
	IDPIntentStateSucceeded
	IDPIntentStateFailed

	idpIntentStateCount
)

func (s IDPIntentState) Valid() bool {
	return s >= 0 && s < idpIntentStateCount
}

func (s IDPIntentState) Exists() bool {
	return s != IDPIntentStateUnspecified && s != IDPIntentStateFailed //TODO: ?
}
