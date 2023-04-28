package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	authn_grpc "github.com/zitadel/zitadel/internal/api/grpc/authn"
	change_grpc "github.com/zitadel/zitadel/internal/api/grpc/change"
	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	project_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/project"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetAppByID(ctx context.Context, req *mgmt_pb.GetAppByIDRequest) (*mgmt_pb.GetAppByIDResponse, error) {
	app, err := s.query.AppByProjectAndAppID(ctx, true, req.ProjectId, req.AppId, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetAppByIDResponse{
		App: project_grpc.AppToPb(app),
	}, nil
}

func (s *Server) ListApps(ctx context.Context, req *mgmt_pb.ListAppsRequest) (*mgmt_pb.ListAppsResponse, error) {
	queries, err := ListAppsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	apps, err := s.query.SearchApps(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppsResponse{
		Result:  project_grpc.AppsToPb(apps.Apps),
		Details: object_grpc.ToListDetails(apps.Count, apps.Sequence, apps.Timestamp),
	}, nil
}

func (s *Server) ListAppChanges(ctx context.Context, req *mgmt_pb.ListAppChangesRequest) (*mgmt_pb.ListAppChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).OrgID).
		AddQuery().
		SequenceGreater(sequence).
		AggregateTypes(project.AggregateType).
		AggregateIDs(req.ProjectId).
		EventData(map[string]interface{}{
			"appId": req.AppId,
		}).
		Builder()
	if asc {
		query.OrderAsc()
	}

	changes, err := s.query.SearchEvents(ctx, query, s.auditLogRetention)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.ListAppChangesResponse{
		Result: change_grpc.EventsToChangesPb(changes, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) AddOIDCApp(ctx context.Context, req *mgmt_pb.AddOIDCAppRequest) (*mgmt_pb.AddOIDCAppResponse, error) {
	appSecretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}
	app, err := s.command.AddOIDCApplication(ctx, AddOIDCAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID, appSecretGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOIDCAppResponse{
		AppId:              app.AppID,
		Details:            object_grpc.AddToDetailsPb(app.Sequence, app.ChangeDate, app.ResourceOwner),
		ClientId:           app.ClientID,
		ClientSecret:       app.ClientSecretString,
		NoneCompliant:      app.Compliance.NoneCompliant,
		ComplianceProblems: project_grpc.ComplianceProblemsToLocalizedMessages(app.Compliance.Problems),
	}, nil
}
func (s *Server) AddSAMLApp(ctx context.Context, req *mgmt_pb.AddSAMLAppRequest) (*mgmt_pb.AddSAMLAppResponse, error) {
	app, err := s.command.AddSAMLApplication(ctx, AddSAMLAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddSAMLAppResponse{
		AppId:   app.AppID,
		Details: object_grpc.AddToDetailsPb(app.Sequence, app.ChangeDate, app.ResourceOwner),
	}, nil
}

func (s *Server) AddAPIApp(ctx context.Context, req *mgmt_pb.AddAPIAppRequest) (*mgmt_pb.AddAPIAppResponse, error) {
	appSecretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}
	app, err := s.command.AddAPIApplication(ctx, AddAPIAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID, appSecretGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddAPIAppResponse{
		AppId:        app.AppID,
		Details:      object_grpc.AddToDetailsPb(app.Sequence, app.ChangeDate, app.ResourceOwner),
		ClientId:     app.ClientID,
		ClientSecret: app.ClientSecretString,
	}, nil
}

func (s *Server) UpdateApp(ctx context.Context, req *mgmt_pb.UpdateAppRequest) (*mgmt_pb.UpdateAppResponse, error) {
	details, err := s.command.ChangeApplication(ctx, req.ProjectId, UpdateAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) UpdateOIDCAppConfig(ctx context.Context, req *mgmt_pb.UpdateOIDCAppConfigRequest) (*mgmt_pb.UpdateOIDCAppConfigResponse, error) {
	config, err := s.command.ChangeOIDCApplication(ctx, UpdateOIDCAppConfigRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOIDCAppConfigResponse{
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateSAMLAppConfig(ctx context.Context, req *mgmt_pb.UpdateSAMLAppConfigRequest) (*mgmt_pb.UpdateSAMLAppConfigResponse, error) {
	config, err := s.command.ChangeSAMLApplication(ctx, UpdateSAMLAppConfigRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateSAMLAppConfigResponse{
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateAPIAppConfig(ctx context.Context, req *mgmt_pb.UpdateAPIAppConfigRequest) (*mgmt_pb.UpdateAPIAppConfigResponse, error) {
	config, err := s.command.ChangeAPIApplication(ctx, UpdateAPIAppConfigRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateAPIAppConfigResponse{
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateApp(ctx context.Context, req *mgmt_pb.DeactivateAppRequest) (*mgmt_pb.DeactivateAppResponse, error) {
	details, err := s.command.DeactivateApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateApp(ctx context.Context, req *mgmt_pb.ReactivateAppRequest) (*mgmt_pb.ReactivateAppResponse, error) {
	details, err := s.command.ReactivateApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RemoveApp(ctx context.Context, req *mgmt_pb.RemoveAppRequest) (*mgmt_pb.RemoveAppResponse, error) {
	details, err := s.command.RemoveApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, req *mgmt_pb.RegenerateOIDCClientSecretRequest) (*mgmt_pb.RegenerateOIDCClientSecretResponse, error) {
	appSecretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}
	config, err := s.command.ChangeOIDCApplicationSecret(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID, appSecretGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RegenerateOIDCClientSecretResponse{
		ClientSecret: config.ClientSecretString,
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) RegenerateAPIClientSecret(ctx context.Context, req *mgmt_pb.RegenerateAPIClientSecretRequest) (*mgmt_pb.RegenerateAPIClientSecretResponse, error) {
	appSecretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}
	config, err := s.command.ChangeAPIApplicationSecret(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID, appSecretGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RegenerateAPIClientSecretResponse{
		ClientSecret: config.ClientSecretString,
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetAppKey(ctx context.Context, req *mgmt_pb.GetAppKeyRequest) (*mgmt_pb.GetAppKeyResponse, error) {
	resourceOwner, err := query.NewAuthNKeyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	aggregateID, err := query.NewAuthNKeyAggregateIDQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	objectID, err := query.NewAuthNKeyObjectIDQuery(req.AppId)
	if err != nil {
		return nil, err
	}
	key, err := s.query.GetAuthNKeyByID(ctx, true, req.KeyId, false, resourceOwner, aggregateID, objectID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetAppKeyResponse{
		Key: authn_grpc.KeyToPb(key),
	}, nil
}

func (s *Server) ListAppKeys(ctx context.Context, req *mgmt_pb.ListAppKeysRequest) (*mgmt_pb.ListAppKeysResponse, error) {
	queries, err := ListAPIClientKeysRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	keys, err := s.query.SearchAuthNKeys(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppKeysResponse{
		Result:  authn_grpc.KeysToPb(keys.AuthNKeys),
		Details: object_grpc.ToListDetails(keys.Count, keys.Sequence, keys.Timestamp),
	}, nil
}

func (s *Server) AddAppKey(ctx context.Context, req *mgmt_pb.AddAppKeyRequest) (*mgmt_pb.AddAppKeyResponse, error) {
	key, err := s.command.AddApplicationKey(ctx, AddAPIClientKeyRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	keyDetails, err := key.Detail()
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddAppKeyResponse{
		Id:         key.KeyID,
		Details:    object_grpc.AddToDetailsPb(key.Sequence, key.ChangeDate, key.ResourceOwner),
		KeyDetails: keyDetails,
	}, nil
}

func (s *Server) RemoveAppKey(ctx context.Context, req *mgmt_pb.RemoveAppKeyRequest) (*mgmt_pb.RemoveAppKeyResponse, error) {
	details, err := s.command.RemoveApplicationKey(ctx, req.ProjectId, req.AppId, req.KeyId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveAppKeyResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}
