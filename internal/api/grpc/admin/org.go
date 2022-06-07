package admin

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	org_grpc "github.com/zitadel/zitadel/internal/api/grpc/org"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	obj_pb "github.com/zitadel/zitadel/pkg/grpc/object"
)

func (s *Server) IsOrgUnique(ctx context.Context, req *admin_pb.IsOrgUniqueRequest) (*admin_pb.IsOrgUniqueResponse, error) {
	isUnique, err := s.query.IsOrgUnique(ctx, req.Name, req.Domain)
	return &admin_pb.IsOrgUniqueResponse{IsUnique: isUnique}, err
}

func (s *Server) GetOrgByID(ctx context.Context, req *admin_pb.GetOrgByIDRequest) (*admin_pb.GetOrgByIDResponse, error) {
	org, err := s.query.OrgByID(ctx, true, req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOrgByIDResponse{Org: org_grpc.OrgViewToPb(org)}, nil
}

func (s *Server) ListOrgs(ctx context.Context, req *admin_pb.ListOrgsRequest) (*admin_pb.ListOrgsResponse, error) {
	queries, err := listOrgRequestToModel(req)
	if err != nil {
		return nil, err
	}
	orgs, err := s.query.SearchOrgs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListOrgsResponse{
		Result: org_grpc.OrgViewsToPb(orgs.Orgs),
		Details: &obj_pb.ListDetails{
			TotalResult:       orgs.Count,
			ProcessedSequence: orgs.Sequence,
			ViewTimestamp:     timestamppb.New(orgs.Timestamp),
		},
	}, nil
}

func (s *Server) SetUpOrg(ctx context.Context, req *admin_pb.SetUpOrgRequest) (*admin_pb.SetUpOrgResponse, error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, domain.NewIAMDomainName(req.Org.Name, s.iamDomain))
	if err != nil {
		return nil, err
	}
	human := setUpOrgHumanToDomain(req.User.(*admin_pb.SetUpOrgRequest_Human_).Human) //TODO: handle machine
	org := setUpOrgOrgToDomain(req.Org)

	objectDetails, err := s.command.SetUpOrg(ctx, org, human, userIDs, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetUpOrgResponse{
		Details: object.DomainToAddDetailsPb(objectDetails),
	}, nil
}

func (s *Server) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgDomain string) ([]string, error) {
	loginName, err := query.NewUserPreferredLoginNameSearchQuery("@"+orgDomain, query.TextEndsWithIgnoreCase)
	if err != nil {
		return nil, err
	}
	users, err := s.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: []query.SearchQuery{loginName}})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, len(users.Users))
	for i, user := range users.Users {
		userIDs[i] = user.ID
	}
	return userIDs, nil
}
