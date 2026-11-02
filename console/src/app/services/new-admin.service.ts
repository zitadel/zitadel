import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  AddEmailProviderHTTPRequestSchema,
  AddEmailProviderSMTPRequestSchema,
  AddSMSProviderHTTPRequestSchema,
  GetDefaultOrgResponse,
  GetMyInstanceResponse,
  SetUpOrgRequestSchema,
  TestEmailProviderSMTPRequestSchema,
  UpdateEmailProviderHTTPRequestSchema,
  UpdateEmailProviderSMTPRequestSchema,
  UpdateSMSProviderHTTPRequestSchema,
} from '@zitadel/proto/zitadel/admin_pb';
import { injectQuery, queryOptions, skipToken } from '@tanstack/angular-query-experimental';
import { NewAuthService } from './new-auth.service';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class NewAdminService {
  constructor(
    private readonly grpcService: GrpcService,
    private readonly authService: NewAuthService,
    private readonly userService: UserService,
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

  public getEmailProviderByIdQueryOptions(id?: string) {
    return queryOptions({
      queryKey: [this.userService.userId(), 'AdminService', 'getEmailProviderById', id],
      queryFn: id ? ({ signal }) => this.getEmailProviderById(id, signal) : skipToken,
    });
  }

  public addEmailProviderSMTP(req: MessageInitShape<typeof AddEmailProviderSMTPRequestSchema>) {
    return this.grpcService.adminNew.addEmailProviderSMTP(req);
  }

  public updateEmailProviderSMTP(req: MessageInitShape<typeof UpdateEmailProviderSMTPRequestSchema>) {
    return this.grpcService.adminNew.updateEmailProviderSMTP(req);
  }

  public addEmailProviderHTTP(req: MessageInitShape<typeof AddEmailProviderHTTPRequestSchema>) {
    return this.grpcService.adminNew.addEmailProviderHTTP(req);
  }

  public updateEmailProviderHTTP(req: MessageInitShape<typeof UpdateEmailProviderHTTPRequestSchema>) {
    return this.grpcService.adminNew.updateEmailProviderHTTP(req);
  }

  public addSMSProviderHTTP(req: MessageInitShape<typeof AddSMSProviderHTTPRequestSchema>) {
    return this.grpcService.adminNew.addSMSProviderHTTP(req);
  }

  public updateSMSProviderHTTP(req: MessageInitShape<typeof UpdateSMSProviderHTTPRequestSchema>) {
    return this.grpcService.adminNew.updateSMSProviderHTTP(req);
  }

  public activateSMTPConfig(id: string) {
    return this.grpcService.adminNew.activateSMTPConfig({ id });
  }

  public deactivateSMTPConfig(id: string) {
    return this.grpcService.adminNew.deactivateSMTPConfig({ id });
  }
}
