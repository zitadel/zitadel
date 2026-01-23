package org

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/user/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func (s *Server) AddOrganization(ctx context.Context, request *connect.Request[org.AddOrganizationRequest]) (*connect.Response[org.AddOrganizationResponse], error) {
	orgSetup, err := addOrganizationRequestToCommand(request.Msg)
	if err != nil {
		return nil, err
	}
	createdOrg, err := s.command.SetUpOrg(ctx, orgSetup, false, s.command.CheckPermissionOrganizationCreate)
	if err != nil {
		return nil, err
	}
	return createdOrganizationToPb(createdOrg)
}

func (s *Server) UpdateOrganization(ctx context.Context, request *connect.Request[org.UpdateOrganizationRequest]) (*connect.Response[org.UpdateOrganizationResponse], error) {
	organization, err := s.command.ChangeOrg(ctx, request.Msg.GetOrganizationId(), request.Msg.GetName(), s.command.CheckPermissionOrganizationWrite)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&org.UpdateOrganizationResponse{
		ChangeDate: timestamppb.New(organization.EventDate),
	}), nil
}

func (s *Server) DeleteOrganization(ctx context.Context, request *connect.Request[org.DeleteOrganizationRequest]) (*connect.Response[org.DeleteOrganizationResponse], error) {
	details, err := s.command.RemoveOrg(ctx, request.Msg.GetOrganizationId(), s.command.CheckPermissionOrganizationDelete, false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.DeleteOrganizationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) SetOrganizationMetadata(ctx context.Context, request *connect.Request[org.SetOrganizationMetadataRequest]) (*connect.Response[org.SetOrganizationMetadataResponse], error) {
	result, err := s.command.BulkSetOrgMetadata(ctx, request.Msg.GetOrganizationId(), s.command.CheckPermissionOrganizationWrite, bulkSetOrgMetadataToDomain(request.Msg)...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.SetOrganizationMetadataResponse{
		SetDate: timestamppb.New(result.EventDate),
	}), nil
}

func (s *Server) DeleteOrganizationMetadata(ctx context.Context, request *connect.Request[org.DeleteOrganizationMetadataRequest]) (*connect.Response[org.DeleteOrganizationMetadataResponse], error) {
	result, err := s.command.BulkRemoveOrgMetadata(ctx, request.Msg.GetOrganizationId(), s.command.CheckPermissionOrganizationWrite, request.Msg.Keys...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.DeleteOrganizationMetadataResponse{
		DeletionDate: timestamppb.New(result.EventDate),
	}), nil
}

func (s *Server) DeactivateOrganization(ctx context.Context, request *connect.Request[org.DeactivateOrganizationRequest]) (*connect.Response[org.DeactivateOrganizationResponse], error) {
	objectDetails, err := s.command.DeactivateOrg(ctx, request.Msg.GetOrganizationId(), s.command.CheckPermissionOrganizationWrite)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.DeactivateOrganizationResponse{
		ChangeDate: timestamppb.New(objectDetails.EventDate),
	}), nil
}

func (s *Server) ActivateOrganization(ctx context.Context, request *connect.Request[org.ActivateOrganizationRequest]) (*connect.Response[org.ActivateOrganizationResponse], error) {
	objectDetails, err := s.command.ReactivateOrg(ctx, request.Msg.GetOrganizationId(), s.command.CheckPermissionOrganizationWrite)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.ActivateOrganizationResponse{
		ChangeDate: timestamppb.New(objectDetails.EventDate),
	}), err
}

func (s *Server) AddOrganizationDomain(ctx context.Context, request *connect.Request[org.AddOrganizationDomainRequest]) (*connect.Response[org.AddOrganizationDomainResponse], error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, request.Msg.GetDomain(), request.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	details, err := s.command.AddOrgDomain(ctx, request.Msg.GetOrganizationId(), request.Msg.GetDomain(), userIDs, s.command.CheckPermissionOrganizationWrite)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.AddOrganizationDomainResponse{
		CreationDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) DeleteOrganizationDomain(ctx context.Context, req *connect.Request[org.DeleteOrganizationDomainRequest]) (*connect.Response[org.DeleteOrganizationDomainResponse], error) {
	details, err := s.command.RemoveOrgDomain(ctx, removeOrgDomainRequestToDomain(req.Msg), s.command.CheckPermissionOrganizationWrite)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.DeleteOrganizationDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), err
}

func (s *Server) GenerateOrganizationDomainValidation(ctx context.Context, req *connect.Request[org.GenerateOrganizationDomainValidationRequest]) (*connect.Response[org.GenerateOrganizationDomainValidationResponse], error) {
	token, url, err := s.command.GenerateOrgDomainValidation(ctx, generateOrgDomainValidationRequestToDomain(req.Msg), s.command.CheckPermissionOrganizationWrite)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.GenerateOrganizationDomainValidationResponse{
		Token: token,
		Url:   url,
	}), nil
}

func (s *Server) VerifyOrganizationDomain(ctx context.Context, request *connect.Request[org.VerifyOrganizationDomainRequest]) (*connect.Response[org.VerifyOrganizationDomainResponse], error) {
	userIDs, err := s.getClaimedUserIDsOfOrgDomain(ctx, request.Msg.GetDomain(), request.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	details, err := s.command.ValidateOrgDomain(ctx, validateOrgDomainRequestToDomain(request.Msg), userIDs, s.command.CheckPermissionOrganizationWrite)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&org.VerifyOrganizationDomainResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}

func addOrganizationRequestToCommand(request *org.AddOrganizationRequest) (*command.OrgSetup, error) {
	admins, err := addOrganizationRequestAdminsToCommand(request.GetAdmins())
	if err != nil {
		return nil, err
	}
	id := request.GetOrganizationId()
	if id == "" {
		id = request.GetOrgId() //nolint:staticcheck
	}
	return &command.OrgSetup{
		Name:         request.GetName(),
		CustomDomain: "",
		Admins:       admins,
		OrgID:        id,
	}, nil
}

func addOrganizationRequestAdminsToCommand(requestAdmins []*org.AddOrganizationRequest_Admin) (admins []*command.OrgSetupAdmin, err error) {
	admins = make([]*command.OrgSetupAdmin, len(requestAdmins))
	for i, admin := range requestAdmins {
		admins[i], err = addOrganizationRequestAdminToCommand(admin)
		if err != nil {
			return nil, err
		}
	}
	return admins, nil
}

func addOrganizationRequestAdminToCommand(admin *org.AddOrganizationRequest_Admin) (*command.OrgSetupAdmin, error) {
	switch a := admin.GetUserType().(type) {
	case *org.AddOrganizationRequest_Admin_UserId:
		return &command.OrgSetupAdmin{
			ID:    a.UserId,
			Roles: admin.GetRoles(),
		}, nil
	case *org.AddOrganizationRequest_Admin_Human:
		human, err := user.AddUserRequestToAddHuman(a.Human)
		if err != nil {
			return nil, err
		}

		return &command.OrgSetupAdmin{
			Human: human,
			Roles: admin.GetRoles(),
		}, nil
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "ORGv2-SD2r1", "userType oneOf %T in method AddOrganization not implemented", a)
	}
}

func createdOrganizationToPb(createdOrg *command.CreatedOrg) (_ *connect.Response[org.AddOrganizationResponse], err error) {
	admins := make([]*org.AddOrganizationResponse_CreatedAdmin, 0, len(createdOrg.OrgAdmins))
	for _, admin := range createdOrg.OrgAdmins {
		admin, ok := admin.(*command.CreatedOrgAdmin)
		if ok {
			admins = append(admins, &org.AddOrganizationResponse_CreatedAdmin{
				UserId:    admin.GetID(),
				EmailCode: admin.EmailCode,
				PhoneCode: admin.PhoneCode,
			})
		}
	}
	return connect.NewResponse(&org.AddOrganizationResponse{
		Details:        object.DomainToDetailsPb(createdOrg.ObjectDetails),
		OrganizationId: createdOrg.ObjectDetails.ResourceOwner,
		CreatedAdmins:  admins,
	}), nil
}

func (s *Server) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgDomain, orgID string) ([]string, error) {
	return s.query.SearchClaimedUserIDsOfOrgDomain(ctx, orgDomain, orgID)
}

func bulkSetOrgMetadataToDomain(req *org.SetOrganizationMetadataRequest) []*domain.Metadata {
	metadata := make([]*domain.Metadata, len(req.Metadata))
	for i, data := range req.Metadata {
		metadata[i] = &domain.Metadata{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metadata
}

func removeOrgDomainRequestToDomain(req *org.DeleteOrganizationDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrganizationId,
		},
		Domain: req.Domain,
	}
}

func generateOrgDomainValidationRequestToDomain(req *org.GenerateOrganizationDomainValidationRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrganizationId,
		},
		Domain:         req.Domain,
		ValidationType: domainValidationTypeToDomain(req.Type),
	}
}

func domainValidationTypeToDomain(validationType org.DomainValidationType) domain.OrgDomainValidationType {
	switch validationType {
	case org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP:
		return domain.OrgDomainValidationTypeHTTP
	case org.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS:
		return domain.OrgDomainValidationTypeDNS
	case org.DomainValidationType_DOMAIN_VALIDATION_TYPE_UNSPECIFIED:
		return domain.OrgDomainValidationTypeUnspecified
	default:
		return domain.OrgDomainValidationTypeUnspecified
	}
}

func validateOrgDomainRequestToDomain(req *org.VerifyOrganizationDomainRequest) *domain.OrgDomain {
	return &domain.OrgDomain{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrganizationId,
		},
		Domain: req.Domain,
	}
}
