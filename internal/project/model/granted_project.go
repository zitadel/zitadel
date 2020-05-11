package model

import (
	"context"
	"github.com/caos/zitadel/internal/api"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/model"
	"time"
)

type GrantedProjectView struct {
	ProjectID       string
	Name            string
	CreationDate    time.Time
	ChangeDate      time.Time
	State           ProjectState
	Type            ProjectType
	ResourceOwner   string
	OrgID           string
	OrgName         string
	OrgDomain       string
	Sequence        uint64
	GrantID         string
	GrantedRoleKeys []string
}

type ProjectType int32

const (
	PROJECTTYPE_OWNED ProjectType = iota
	PROJECTTYPE_GRANTED
)

type GrantedProjectSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn GrantedProjectSearchKey
	Asc           bool
	Queries       []*GrantedProjectSearchQuery
}

type GrantedProjectSearchKey int32

const (
	GRANTEDPROJECTSEARCHKEY_UNSPECIFIED GrantedProjectSearchKey = iota
	GRANTEDPROJECTSEARCHKEY_NAME
	GRANTEDPROJECTSEARCHKEY_PROJECTID
	GRANTEDPROJECTSEARCHKEY_GRANTID
	GRANTEDPROJECTSEARCHKEY_ORGID
	GRANTEDPROJECTSEARCHKEY_RESOURCE_OWNER
)

type GrantedProjectSearchQuery struct {
	Key    GrantedProjectSearchKey
	Method model.SearchMethod
	Value  string
}

type GrantedProjectSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*GrantedProjectView
}

func (r *GrantedProjectSearchRequest) AppendMyOrgQuery(ctx context.Context) {
	orgID := grpc_util.GetHeader(ctx, api.ZitadelOrgID)
	r.Queries = append(r.Queries, &GrantedProjectSearchQuery{Key: GRANTEDPROJECTSEARCHKEY_ORGID, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}

func (r *GrantedProjectSearchRequest) AppendNotMyOrgQuery(ctx context.Context) {
	orgID := grpc_util.GetHeader(ctx, api.ZitadelOrgID)
	r.Queries = append(r.Queries, &GrantedProjectSearchQuery{Key: GRANTEDPROJECTSEARCHKEY_ORGID, Method: model.SEARCHMETHOD_NOT_EQUALS, Value: orgID})
}

func (r *GrantedProjectSearchRequest) AppendMyResourceOwnerQuery(ctx context.Context) {
	orgID := grpc_util.GetHeader(ctx, api.ZitadelOrgID)
	r.Queries = append(r.Queries, &GrantedProjectSearchQuery{Key: GRANTEDPROJECTSEARCHKEY_RESOURCE_OWNER, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}

func (r *GrantedProjectSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
