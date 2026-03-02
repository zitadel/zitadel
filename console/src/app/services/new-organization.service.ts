import { computed, Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { injectQuery, mutationOptions, QueryClient, queryOptions, skipToken } from '@tanstack/angular-query-experimental';
import { create, DescMessage, MessageInitShape, toBinary } from '@bufbuild/protobuf';
import {
  ListOrganizationsRequestSchema,
  ListOrganizationsResponse,
  OrganizationService,
} from '@zitadel/proto/zitadel/org/v2/org_service_pb';
import { NewMgmtService } from './new-mgmt.service';
import { OrgInterceptorProvider } from './interceptors/org.interceptor';
import { NewAdminService } from './new-admin.service';
import { SetUpOrgRequestSchema } from '@zitadel/proto/zitadel/admin_pb';
import { UserService } from './user.service';
import { GrpcAuthService } from './grpc-auth.service';
import { concatWith, defer, map } from 'rxjs';
import { filter } from 'rxjs/operators';
import { toSignal } from '@angular/core/rxjs-interop';
import { Buffer } from 'buffer';
import { AuthService, ListMyProjectOrgsRequestSchema } from '@zitadel/proto/zitadel/auth_pb';

@Injectable({
  providedIn: 'root',
})
export class NewOrganizationService {
  constructor(
    private readonly grpcService: GrpcService,
    private readonly authService: GrpcAuthService,
    private readonly newMgtmService: NewMgmtService,
    private readonly newAdminService: NewAdminService,
    private readonly orgInterceptorProvider: OrgInterceptorProvider,
    private readonly queryClient: QueryClient,
    private readonly userService: UserService,
  ) {}

  public organizationByIdQueryOptions(organizationId?: string) {
    const req = {
      query: {
        limit: 1,
      },
      queries: [
        {
          query: {
            case: 'idQuery' as const,
            value: {
              id: organizationId?.toString(),
            },
          },
        },
      ],
    };

    const { queryFn, ...listOrganizationsQueryOptions } = this.listOrganizationsQueryOptions(req);

    return queryOptions({
      ...listOrganizationsQueryOptions,
      queryFn: organizationId ? queryFn : skipToken,
      select: (data) => data.result.find(Boolean) ?? null,
    });
  }

  public activeOrganizationQuery() {
    const activeOrg$ = defer(() => this.authService.getActiveOrg()).pipe(
      concatWith(this.authService.activeOrgChanged),
      filter(Boolean),
      map((org) => org.id),
    );

    const activeOrg = toSignal(activeOrg$);

    const req = computed(
      () =>
        ({
          query: {
            limit: 1,
          },
          queries: [
            {
              query: {
                case: 'idQuery' as const,
                value: {
                  id: activeOrg(),
                },
              },
            },
          ],
        }) satisfies MessageInitShape<typeof ListMyProjectOrgsRequestSchema>,
    );

    return injectQuery(() => {
      const { queryFn, ...listMyProjectOrgsQueryOptions } = this.listMyProjectOrgsQueryOptions(req());

      return queryOptions({
        ...listMyProjectOrgsQueryOptions,
        queryFn: activeOrg() ? queryFn : skipToken,
        select: (data) => data.result.find(Boolean) ?? null,
      });
    });
  }

  public listOrganizationsQueryOptions(req?: MessageInitShape<typeof ListOrganizationsRequestSchema>) {
    const queryKeyHashFn = tanstackQueryKeyHashFn(ListOrganizationsRequestSchema);
    return queryOptions({
      queryKey: [
        this.userService.userId(),
        OrganizationService.name,
        OrganizationService.method.listOrganizations.name,
        req,
      ] as const,
      queryKeyHashFn: (key) => queryKeyHashFn(...key),
      queryFn: ({ signal }) => this.listOrganizations(req ?? {}, signal),
    });
  }

  private listOrganizations(
    req: MessageInitShape<typeof ListOrganizationsRequestSchema>,
    signal?: AbortSignal,
  ): Promise<ListOrganizationsResponse> {
    return this.grpcService.organizationNew.listOrganizations(req, { signal });
  }

  private listMyProjectOrgsQueryOptions(req?: MessageInitShape<typeof ListMyProjectOrgsRequestSchema>) {
    const queryKeyHashFn = tanstackQueryKeyHashFn(ListMyProjectOrgsRequestSchema);
    return queryOptions({
      queryKey: [this.userService.userId(), AuthService.name, AuthService.method.listMyProjectOrgs.name, req] as const,
      queryKeyHashFn: (key) => queryKeyHashFn(...key),
      queryFn: ({ signal }) => this.grpcService.authNew.listMyProjectOrgs(req ?? {}, { signal }),
    });
  }

  private invalidateAllOrganizationQueries() {
    return this.queryClient.invalidateQueries({
      queryKey: this.listOrganizationsQueryOptions().queryKey,
    });
  }

  public renameOrgMutationOptions = () =>
    mutationOptions({
      mutationKey: ['renameOrg'],
      mutationFn: (name: string) => this.newMgtmService.updateOrg({ name }),
      onSettled: () => this.invalidateAllOrganizationQueries(),
    });

  public deleteOrgMutationOptions = () =>
    mutationOptions({
      mutationKey: ['deleteOrg'],
      mutationFn: async () => {
        const resp = await this.newMgtmService.removeOrg();

        // We change active org to default org as
        // current org was deleted to avoid Organization doesn't exist
        await this.authService.getActiveOrg();

        return resp;
      },
      onSettled: async () => {
        const orgId = this.orgInterceptorProvider.getOrgId();
        if (orgId) {
          this.queryClient.removeQueries({
            queryKey: this.organizationByIdQueryOptions(orgId).queryKey,
          });
        }

        await this.invalidateAllOrganizationQueries();
      },
    });

  public reactivateOrgMutationOptions = () =>
    mutationOptions({
      mutationKey: ['reactivateOrg'],
      mutationFn: () => this.newMgtmService.reactivateOrg(),
      onSettled: () => this.invalidateAllOrganizationQueries(),
    });

  public deactivateOrgMutationOptions = () =>
    mutationOptions({
      mutationKey: ['deactivateOrg'],
      mutationFn: () => this.newMgtmService.deactivateOrg(),
      onSettled: () => this.invalidateAllOrganizationQueries(),
    });

  public setupOrgMutationOptions = () =>
    mutationOptions({
      mutationKey: ['setupOrg'],
      mutationFn: (req: MessageInitShape<typeof SetUpOrgRequestSchema>) => this.newAdminService.setupOrg(req),
      onSettled: async () => {
        await this.invalidateAllOrganizationQueries();
      },
    });

  public addOrgMutationOptions = () =>
    mutationOptions({
      mutationKey: ['addOrg'],
      mutationFn: (name: string) => this.newMgtmService.addOrg(name),
      onSettled: async () => {
        await this.invalidateAllOrganizationQueries();
      },
    });
}

function tanstackQueryKeyHashFn<T extends DescMessage>(schema: T) {
  return (userId: string | undefined, serviceName: string, methodName: string, req?: MessageInitShape<T>) => {
    if (!req) {
      return JSON.stringify([userId, serviceName, methodName]);
    }

    const serializedReq = toBinary(schema, create(schema, req));
    const serializedReqAsString = Buffer.from(serializedReq).toString('base64');

    return JSON.stringify([userId, serviceName, methodName, serializedReqAsString]);
  };
}
