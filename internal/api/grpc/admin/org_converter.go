package admin

import (
	"github.com/caos/logging"
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/org/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func listOrgRequestToModel(req *admin.ListOrgsRequest) (*model.OrgSearchRequest, error) {
	queries, err := org_grpc.OrgQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &model.OrgSearchRequest{
		Offset:  req.MetaData.Offset,
		Limit:   uint64(req.MetaData.Limit),
		Asc:     req.MetaData.Asc,
		Queries: queries,
	}, nil
}

func setUpOrgOrgToDomain(req *admin.SetUpOrgRequest_Org) *domain.Org {
	org := &domain.Org{
		Name:    req.Name,
		Domains: []*domain.OrgDomain{},
	}
	if req.Domain != "" {
		org.Domains = append(org.Domains, &domain.OrgDomain{Domain: req.Domain})
	}
	return org
}
