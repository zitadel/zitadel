package project

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func (s *Server) CreateProjectGrant(ctx context.Context, req *project_pb.CreateProjectGrantRequest) (*project_pb.CreateProjectGrantResponse, error) {
	add := projectGrantCreateToCommand(req)
	project, err := s.command.AddProjectGrant(ctx, add)
	if err != nil {
		return nil, err
	}
	var creationDate *timestamppb.Timestamp
	if !project.EventDate.IsZero() {
		creationDate = timestamppb.New(project.EventDate)
	}
	return &project_pb.CreateProjectGrantResponse{
		CreationDate: creationDate,
	}, nil
}

func projectGrantCreateToCommand(req *project_pb.CreateProjectGrantRequest) *command.AddProjectGrant {
	return &command.AddProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		GrantedOrgID: req.GrantedOrganizationId,
		RoleKeys:     req.RoleKeys,
	}
}
