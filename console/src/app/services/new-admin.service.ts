import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  AddEmailProviderSMTPRequestSchema,
  GetDefaultOrgResponse,
  GetMyInstanceResponse,
  SetUpOrgRequestSchema,
  TestEmailProviderSMTPRequestSchema,
  UpdateEmailProviderSMTPRequestSchema,
} from '@zitadel/proto/zitadel/admin_pb';
import { injectQuery, QueryClient, queryOptions, skipToken } from '@tanstack/angular-query-experimental';
import { NewAuthService } from './new-auth.service';
import { UserService } from './user.service';
import { NewSettingsService } from './new-settings.service';

@Injectable({
  providedIn: 'root',
})
export class NewAdminService {
  constructor(
    private readonly grpcService: GrpcService,
    private readonly authService: NewAuthService,
    private readonly userService: UserService,
    private readonly queryClient: QueryClient,
    private readonly settingsService: NewSettingsService,
  ) {}

  public setupOrg(req: MessageInitShape<typeof SetUpOrgRequestSchema>) {
    return this.grpcService.adminNew.setUpOrg(req);
  }

  public getDefaultOrg(): Promise<GetDefaultOrgResponse> {
    return this.grpcService.adminNew.getDefaultOrg({});
  }

  private getMyInstance(signal?: AbortSignal): Promise<GetMyInstanceResponse> {
    return this.grpcService.adminNew.getMyInstance({}, { signal });
  }

  public getMyInstanceQuery() {
    const listMyZitadelPermissionsQuery = this.authService.listMyZitadelPermissionsQuery();
    return injectQuery(() => ({
      queryKey: [this.userService.userId(), 'admin', 'getMyInstance'],
      queryFn: async () => this.getMyInstance(),
      enabled: (listMyZitadelPermissionsQuery.data() ?? []).includes('iam.write'),
    }));
  }

  public testEmailProviderSMTP(req: MessageInitShape<typeof TestEmailProviderSMTPRequestSchema>) {
    return this.grpcService.adminNew.testEmailProviderSMTP(req);
  }

  public getEmailProviderById(id: string, signal: AbortSignal) {
    return this.grpcService.adminNew.getEmailProviderById({ id }, { signal });
  }

  public listEmailProvidersQueryOptions() {
    return queryOptions({
      queryKey: [this.userService.userId(), 'AdminService', 'emailProviders', 'list'],
      queryFn: ({ signal }) => this.grpcService.adminNew.listEmailProviders({}, { signal }),
    });
  }

  // Reloads the email providers after one of them was activated, deactivated, changed or removed.
  // The general settings state the restrictions of the active provider, so they are reloaded as well.
  public invalidateEmailProviders() {
    return Promise.all([
      this.queryClient.invalidateQueries({
        queryKey: [this.userService.userId(), 'AdminService', 'emailProviders'],
      }),
      this.settingsService.invalidateGeneralSettings(),
    ]);
  }

  public getEmailProviderByIdQueryOptions(id?: string) {
    return queryOptions({
      queryKey: [this.userService.userId(), 'AdminService', 'getEmailProviderById', id],
      queryFn: id ? ({ signal }) => this.getEmailProviderById(id, signal) : skipToken,
    });
  }

  public addEmailProviderSMTP(req: MessageInitShape<typeof AddEmailProviderSMTPRequestSchema>) {
    return this.grpcService.adminNew.addEmailProviderSMTP(req).finally(() => this.invalidateEmailProviders());
  }

  public updateEmailProviderSMTP(req: MessageInitShape<typeof UpdateEmailProviderSMTPRequestSchema>) {
    return this.grpcService.adminNew.updateEmailProviderSMTP(req).finally(() => this.invalidateEmailProviders());
  }

  public activateSMTPConfig(id: string) {
    return this.grpcService.adminNew.activateSMTPConfig({ id }).finally(() => this.invalidateEmailProviders());
  }

  public deactivateSMTPConfig(id: string) {
    return this.grpcService.adminNew.deactivateSMTPConfig({ id }).finally(() => this.invalidateEmailProviders());
  }
}
