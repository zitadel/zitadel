package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	org_grpc "github.com/zitadel/zitadel/internal/api/grpc/org"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) IsOrgUnique(ctx context.Context, req *admin_pb.IsOrgUniqueRequest) (*admin_pb.IsOrgUniqueResponse, error) {
	isUnique, err := s.query.IsOrgUnique(ctx, req.Name, req.Domain)
	return &admin_pb.IsOrgUniqueResponse{IsUnique: isUnique}, err
}

func (s *Server) SetDefaultOrg(ctx context.Context, req *admin_pb.SetDefaultOrgRequest) (*admin_pb.SetDefaultOrgResponse, error) {
	details, err := s.command.SetDefaultOrg(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultOrgResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RemoveOrg(ctx context.Context, req *admin_pb.RemoveOrgRequest) (*admin_pb.RemoveOrgResponse, error) {
	details, err := s.command.RemoveOrg(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveOrgResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) GetDefaultOrg(ctx context.Context, _ *admin_pb.GetDefaultOrgRequest) (*admin_pb.GetDefaultOrgResponse, error) {
	org, err := s.query.OrgByID(ctx, true, authz.GetInstance(ctx).DefaultOrganisationID())
	return &admin_pb.GetDefaultOrgResponse{Org: org_grpc.OrgToPb(org)}, err
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
		Result:  org_grpc.OrgViewsToPb(orgs.Orgs),
		Details: object.ToListDetails(orgs.Count, orgs.Sequence, orgs.LastRun),
	}, nil
}

func (s *Server) SetUpOrg(ctx context.Context, req *admin_pb.SetUpOrgRequest) (*admin_pb.SetUpOrgResponse, error) {
	orgDomain, err := domain.NewIAMDomainName(req.Org.Name, authz.GetInstance(ctx).RequestedDomain())
	if err != nil {
		return nil, err
	}
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, orgDomain)
	if err != nil {
		return nil, err
	}
	human := setUpOrgHumanToCommand(req.User.(*admin_pb.SetUpOrgRequest_Human_).Human) //TODO: handle machine

	createdOrg, err := s.command.SetUpOrg(ctx, &command.OrgSetup{
		Name:         req.Org.Name,
		CustomDomain: req.Org.Domain,
		Admins: []*command.OrgSetupAdmin{
			{
				Human: human,
				Roles: req.Roles,
			},
		},
	}, true, userIDs...)
	if err != nil {
		return nil, err
	}
	var userID string
	if len(createdOrg.CreatedAdmins) == 1 {
		userID = createdOrg.CreatedAdmins[0].ID
	}
	return &admin_pb.SetUpOrgResponse{
		Details: object.DomainToAddDetailsPb(createdOrg.ObjectDetails),
		OrgId:   createdOrg.ObjectDetails.ResourceOwner,
		UserId:  userID,
	}, nil
}

func (s *Server) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgDomain string) ([]string, error) {
	loginName, err := query.NewUserPreferredLoginNameSearchQuery("@"+orgDomain, query.TextEndsWithIgnoreCase)
	if err != nil {
		return nil, err
	}
	users, err := s.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: []query.SearchQuery{loginName}}, false)
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
