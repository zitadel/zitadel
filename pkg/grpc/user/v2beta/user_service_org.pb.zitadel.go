package user

import "github.com/zitadel/zitadel/internal/api/grpc/server/middleware"

// OrganisationFromRequest implements deprecated [middleware.OrganisationFromRequest] interface.
// it will be removed before going GA (https://github.com/zitadel/zitadel/issues/6718)
func (r *AddHumanUserRequest) OrganisationFromRequest() *middleware.Organization {
	return &middleware.Organization{
		ID:     r.GetOrganisation().GetOrgId(),
		Domain: r.GetOrganisation().GetOrgDomain(),
	}
}
