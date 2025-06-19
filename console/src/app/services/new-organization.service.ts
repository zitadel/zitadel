import { computed, Injectable, signal } from '@angular/core';
import { GrpcService } from './grpc.service';
import { injectQuery, mutationOptions, QueryClient, queryOptions, skipToken } from '@tanstack/angular-query-experimental';
import { MessageInitShape } from '@bufbuild/protobuf';
import { ListOrganizationsRequestSchema, ListOrganizationsResponse } from '@zitadel/proto/zitadel/org/v2/org_service_pb';
import { NewMgmtService } from './new-mgmt.service';
import { OrgInterceptorProvider } from './interceptors/org.interceptor';
import { NewAdminService } from './new-admin.service';
import { SetUpOrgRequestSchema } from '@zitadel/proto/zitadel/admin_pb';
import { TranslateService } from '@ngx-translate/core';
import { lastValueFrom } from 'rxjs';
import { first } from 'rxjs/operators';
import { StorageKey, StorageLocation, StorageService } from './storage.service';

@Injectable({
  providedIn: 'root',
})
export class NewOrganizationService {
  constructor(
    private readonly grpcService: GrpcService,
    private readonly newMgtmService: NewMgmtService,
    private readonly newAdminService: NewAdminService,
    private readonly orgInterceptorProvider: OrgInterceptorProvider,
    private readonly queryClient: QueryClient,
    private readonly translate: TranslateService,
    private readonly storage: StorageService,
  ) {}

  private readonly orgIdSignal = signal<string | undefined>(
    this.storage.getItem(StorageKey.organizationId, StorageLocation.session) ??
      this.storage.getItem(StorageKey.organizationId, StorageLocation.local) ??
      undefined,
  );
  public readonly orgId = this.orgIdSignal.asReadonly();

  public getOrgId() {
    return computed(() => {
      const orgId = this.orgIdSignal();
      if (orgId === undefined) {
        throw new Error('No organization ID set');
      }
      return orgId;
    });
  }

  public async setOrgId(orgId?: string) {
    console.log('beboop', orgId);
    console.trace(orgId);
    const organization = await this.queryClient.fetchQuery(this.organizationByIdQueryOptions(orgId ?? this.getOrgId()()));
    if (organization) {
      this.storage.setItem(StorageKey.organizationId, orgId, StorageLocation.session);
      this.storage.setItem(StorageKey.organizationId, orgId, StorageLocation.local);
      this.orgIdSignal.set(orgId);
    } else {
      throw new Error('request organization not found');
    }
    return organization;
  }

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

    return queryOptions({
      queryKey: ['organization', 'listOrganizations', req],
      queryFn: organizationId
        ? () => this.listOrganizations(req).then((resp) => resp.result.find(Boolean) ?? null)
        : skipToken,
    });
  }

  public activeOrganizationQuery() {
    return injectQuery(() => this.organizationByIdQueryOptions(this.orgId()));
  }

  public listOrganizationsQueryOptions(req?: MessageInitShape<typeof ListOrganizationsRequestSchema>) {
    return queryOptions({
      queryKey: this.listOrganizationsQueryKey(req),
      queryFn: () => this.listOrganizations(req ?? {}),
    });
  }

  public listOrganizationsQueryKey(req?: MessageInitShape<typeof ListOrganizationsRequestSchema>) {
    if (!req) {
      return ['organization', 'listOrganizations'];
    }

    // needed because angular query isn't able to serialize a bigint key
    const query = req.query ? { ...req.query, offset: req.query.offset ? Number(req.query.offset) : undefined } : undefined;
    const queryKey = {
      ...req,
      ...(query ? { query } : {}),
    };

    return ['organization', 'listOrganizations', queryKey];
  }

  public listOrganizations(
    req: MessageInitShape<typeof ListOrganizationsRequestSchema>,
    signal?: AbortSignal,
  ): Promise<ListOrganizationsResponse> {
    return this.grpcService.organizationNew.listOrganizations(req, { signal });
  }

  private async getDefaultOrganization() {
    let resp = await this.listOrganizations({
      query: {
        limit: 1,
      },
      queries: [
        {
          query: {
            case: 'defaultQuery',
            value: {},
          },
        },
      ],
    });
    return resp.result.find(Boolean) ?? null;
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
        // Before we remove the org we get the current default org
        // we have to query before the current org is removed
        const defaultOrg = await this.getDefaultOrganization();
        if (!defaultOrg) {
          const error$ = this.translate.get('ORG.TOAST.DEFAULTORGOTFOUND').pipe(first());
          throw { message: await lastValueFrom(error$) };
        }

        const resp = await this.newMgtmService.removeOrg();
        await new Promise((resolve) => setTimeout(resolve, 1000));

        // We change active org to default org as
        // current org was deleted to avoid Organization doesn't exist
        await this.setOrgId(defaultOrg.id);

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
