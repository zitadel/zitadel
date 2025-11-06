package instancev2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api/instance/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/api/authz"
	filter_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	filter_v2beta "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	instance_v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

// =================
// v2Beta endpoints
// =================

func AddCustomDomainBeta(ctx context.Context, request *connect.Request[instance_v2beta.AddCustomDomainRequest]) (*connect.Response[instance_v2beta.AddCustomDomainResponse], error) {
	addCustomDomainCmd := domain.NewAddInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetDomain(), domain.DomainTypeCustom)
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

func ListCustomDomainsBeta(ctx context.Context, request *connect.Request[instance_v2beta.ListCustomDomainsRequest]) (*connect.Response[instance_v2beta.ListCustomDomainsResponse], error) {
	listCustomDomainsQuery := domain.NewListInstanceDomainsQuery(convert.ListCustomDomainsBetaRequestToV2Request(request.Msg))

	err := domain.Invoke(ctx, listCustomDomainsQuery, domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()))
	if err != nil {
		return nil, err
	}

	customDomains := listCustomDomainsQuery.Result()
	return &connect.Response[instance_v2beta.ListCustomDomainsResponse]{
		Msg: &instance_v2beta.ListCustomDomainsResponse{
			Domains: convert.DomainInstanceDomainListModelToGRPCBetaResponse(customDomains),
			Pagination: &filter_v2beta.PaginationResponse{
				// TODO(IAM-Marco): return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10955
				TotalResult:  uint64(len(customDomains)),
				AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
			},
		},
	}, nil
}

func AddTrustedDomainBeta(ctx context.Context, request *connect.Request[instance_v2beta.AddTrustedDomainRequest]) (*connect.Response[instance_v2beta.AddTrustedDomainResponse], error) {
	addTrustedDomainCmd := domain.NewAddInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetDomain(), domain.DomainTypeTrusted)

	err := domain.Invoke(
		ctx,
		addTrustedDomainCmd,
		domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()),
	)

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2beta.AddTrustedDomainResponse]{
		Msg: &instance_v2beta.AddTrustedDomainResponse{
			// TODO(IAM-Marco): Return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10881
			CreationDate: timestamppb.Now(),
		},
	}, nil
}

// =================
// v2 endpoints
// =================

func AddCustomDomain(ctx context.Context, request *connect.Request[instance_v2.AddCustomDomainRequest]) (*connect.Response[instance_v2.AddCustomDomainResponse], error) {
	addCustomDomainCmd := domain.NewAddInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetCustomDomain(), domain.DomainTypeCustom)
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

func ListCustomDomains(ctx context.Context, request *connect.Request[instance_v2.ListCustomDomainsRequest]) (*connect.Response[instance_v2.ListCustomDomainsResponse], error) {
	listCustomDomainsQuery := domain.NewListInstanceDomainsQuery(request.Msg)

	err := domain.Invoke(ctx, listCustomDomainsQuery, domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()))
	if err != nil {
		return nil, err
	}

	customDomains := listCustomDomainsQuery.Result()
	return &connect.Response[instance_v2.ListCustomDomainsResponse]{
		Msg: &instance_v2.ListCustomDomainsResponse{
			Domains: convert.DomainInstanceDomainListModelToGRPCResponse(customDomains),
			Pagination: &filter_v2.PaginationResponse{
				// TODO(IAM-Marco): return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10955
				TotalResult:  uint64(len(customDomains)),
				AppliedLimit: uint64(request.Msg.GetPagination().GetLimit()),
			},
		},
	}, nil
}

func AddTrustedDomain(ctx context.Context, request *connect.Request[instance_v2.AddTrustedDomainRequest]) (*connect.Response[instance_v2.AddTrustedDomainResponse], error) {
	addTrustedDomainCmd := domain.NewAddInstanceDomainCommand(request.Msg.GetInstanceId(), request.Msg.GetTrustedDomain(), domain.DomainTypeTrusted)

	err := domain.Invoke(
		ctx,
		addTrustedDomainCmd,
		domain.WithInstanceDomainRepo(repository.InstanceDomainRepository()),
	)

	if err != nil {
		return nil, err
	}

	return &connect.Response[instance_v2.AddTrustedDomainResponse]{
		Msg: &instance_v2.AddTrustedDomainResponse{
			// TODO(IAM-Marco): Return correct value. Tracked in https://github.com/zitadel/zitadel/issues/10881
			CreationDate: timestamppb.Now(),
		},
	}, nil
}
