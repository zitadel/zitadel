package admin

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	org_grpc "github.com/zitadel/zitadel/internal/api/grpc/org"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	cmd_v2 "github.com/zitadel/zitadel/internal/v2/command"
	"github.com/zitadel/zitadel/internal/v2/org"
	"github.com/zitadel/zitadel/internal/v2/projection"
	"github.com/zitadel/zitadel/internal/v2/readmodel"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
	org_pb "github.com/zitadel/zitadel/pkg/grpc/org"
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
	intent, err := cmd_v2.NewRemoveOrg(req.GetOrgId()).ToPushIntent(ctx, s.es.Querier)
	if err != nil {
		return nil, err
	}
	if intent == nil {
		return new(admin_pb.RemoveOrgResponse), nil
	}
	return new(admin_pb.RemoveOrgResponse), s.es.Push(ctx, intent)
}

func (s *Server) GetDefaultOrg(ctx context.Context, _ *admin_pb.GetDefaultOrgRequest) (*admin_pb.GetDefaultOrgResponse, error) {
	// org, err := s.query.OrgByID(ctx, true, authz.GetInstance(ctx).DefaultOrganisationID())
	// return &admin_pb.GetDefaultOrgResponse{Org: org_grpc.OrgToPb(org)}, err

	org := readmodel.NewOrg(authz.GetInstance(ctx).DefaultOrganisationID())

	// if err := s.es.Query(ctx, authz.GetInstance(ctx).InstanceID(), org, org.Filter()...); err != nil {
	// 	return nil, err
	// }

	return &admin_pb.GetDefaultOrgResponse{Org: orgToPb(org)}, nil
}

func (s *Server) GetOrgByID(ctx context.Context, req *admin_pb.GetOrgByIDRequest) (*admin_pb.GetOrgByIDResponse, error) {
	org := readmodel.NewOrg(req.GetId())

	// if err := s.es.Query(ctx, authz.GetInstance(ctx).InstanceID(), org, org.Filter()...); err != nil {
	// 	return nil, err
	// }

	return &admin_pb.GetOrgByIDResponse{Org: orgToPb(org)}, nil
	// org, err := s.query.OrgByID(ctx, true, req.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// return &admin_pb.GetOrgByIDResponse{Org: org_grpc.OrgViewToPb(org)}, nil
}

func orgToPb(org *readmodel.Org) *org_pb.Org {
	res := &org_pb.Org{
		Id:            org.ID,
		State:         stateToPb(org.State),
		Name:          org.Name,
		PrimaryDomain: org.PrimaryDomain.Domain,
		Details: &object_pb.ObjectDetails{
			Sequence:      uint64(org.Sequence),
			CreationDate:  timestamppb.New(org.CreationDate),
			ChangeDate:    timestamppb.New(org.ChangeDate),
			ResourceOwner: org.Owner,
		},
	}

	if !org.CreationDate.IsZero() {
		res.Details.CreationDate = timestamppb.New(org.CreationDate)
	}

	if !org.ChangeDate.IsZero() {
		res.Details.ChangeDate = timestamppb.New(org.ChangeDate)
	}

	return res
}

func stateToPb(state *projection.OrgState) org_pb.OrgState {
	switch state.State {
	case org.ActiveState:
		return org_pb.OrgState_ORG_STATE_ACTIVE
	case org.InactiveState:
		return org_pb.OrgState_ORG_STATE_INACTIVE
	case org.RemovedState:
		return org_pb.OrgState_ORG_STATE_REMOVED
	default:
		return org_pb.OrgState_ORG_STATE_UNSPECIFIED
	}
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
