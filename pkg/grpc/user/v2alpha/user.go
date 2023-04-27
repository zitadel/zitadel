package user

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func (r *AddHumanUserRequest) AuthContext() string {
	return r.GetOrganisation().GetOrgId()
}

func (r *AddHumanUserResponse) Location() string {
	return fmt.Sprintf("%s/api/%s", http_util.BuildOrigin(authz.GetInstance(context.TODO()).RequestedHost(), false), r.UserId)
}
