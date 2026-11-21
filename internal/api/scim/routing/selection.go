// Package routing contains conservative policy checks for SCIM user login routing.
// Routing to an IdP is not proof of ownership of an external identity.
package routing

// Decision describes whether an explicitly configured provider can be offered
// for a SCIM-provisioned user's first authentication.
type Decision struct {
	ProviderID string
	Allowed    bool
}

// Select returns an IdP only when the trusted provisioning configuration
// names exactly one provider that is active and enabled for the organization.
// No subject linking is performed here. Never derive candidates from SCIM
// externalId, email address, username, or other caller-supplied attributes.
func Select(configuredProviderIDs []string, activeProviderIDs []string, externalLoginAllowed bool) Decision {
	if !externalLoginAllowed || len(configuredProviderIDs) != 1 || configuredProviderIDs[0] == "" {
		return Decision{}
	}
	for _, active := range activeProviderIDs {
		if active == configuredProviderIDs[0] {
			return Decision{ProviderID: active, Allowed: true}
		}
	}
	return Decision{}
}
