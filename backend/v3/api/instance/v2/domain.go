package instancev2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/api/authz"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	instance_v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

// =================
// v2Beta endpoints
// =================

func AddCustomDomainBeta(ctx context.Context, request *connect.Request[instance_v2beta.AddCustomDomainRequest]) (*connect.Response[instance_v2beta.AddCustomDomainResponse], error) {
	addCustomDomainCmd := domain.NewAddInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetDomain())
	oidcConfigUpdateCmd := domain.NewOIDCConfigurationUpdate(request.Msg.GetDomain(), authz.GetInstance(ctx).ProjectID(), authz.GetInstance(ctx).ConsoleApplicationID())

	batchExec := domain.BatchExecutors(
		addCustomDomainCmd,
		oidcConfigUpdateCmd,
	)

	err := domain.Invoke(
		ctx,
		batchExec,
		domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()),
		// domain.WithOIDCConfigurationRepo(repository.OIDCConfigurationRepository()),
	)

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2beta.AddCustomDomainResponse]{
		Msg: &instance_v2beta.AddCustomDomainResponse{
			// TODO(IAM-Marco): Return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10881
			CreationDate: timestamppb.Now(),
		},
	}, nil
}

func RemoveCustomDomainBeta(ctx context.Context, request *connect.Request[instance_v2beta.RemoveCustomDomainRequest]) (*connect.Response[instance_v2beta.RemoveCustomDomainResponse], error) {
	removeCustomDomainCmd := domain.NewRemoveInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetDomain())

	err := domain.Invoke(ctx, removeCustomDomainCmd, domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()))
	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2beta.RemoveCustomDomainResponse]{
		Msg: &instance_v2beta.RemoveCustomDomainResponse{
			// TODO(IAM-Marco): Return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10881
			DeletionDate: timestamppb.Now(),
		},
	}, nil
}

// =================
// v2 endpoints
// =================

func AddCustomDomain(ctx context.Context, request *connect.Request[instance_v2.AddCustomDomainRequest]) (*connect.Response[instance_v2.AddCustomDomainResponse], error) {
	addCustomDomainCmd := domain.NewAddInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetCustomDomain())
	oidcConfigUpdateCmd := domain.NewOIDCConfigurationUpdate(request.Msg.GetCustomDomain(), authz.GetInstance(ctx).ProjectID(), authz.GetInstance(ctx).ConsoleApplicationID())

	batchExec := domain.BatchExecutors(
		addCustomDomainCmd,
		oidcConfigUpdateCmd,
	)

	err := domain.Invoke(
		ctx,
		batchExec,
		domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()),
		// domain.WithOIDCConfigurationRepo(repository.OIDCConfigurationRepository()),
	)

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2.AddCustomDomainResponse]{
		Msg: &instance_v2.AddCustomDomainResponse{
			// TODO(IAM-Marco): Return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10881
			CreationDate: timestamppb.Now(),
		},
	}, nil
}

func RemoveCustomDomain(ctx context.Context, request *connect.Request[instance_v2.RemoveCustomDomainRequest]) (*connect.Response[instance_v2.RemoveCustomDomainResponse], error) {
	removeCustomDomainCmd := domain.NewRemoveInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetCustomDomain())

	err := domain.Invoke(ctx, removeCustomDomainCmd, domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()))
	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2.RemoveCustomDomainResponse]{
		Msg: &instance_v2.RemoveCustomDomainResponse{
			// TODO(IAM-Marco): Return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10881
			DeletionDate: timestamppb.Now(),
		},
	}, nil
}
