import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  GetDefaultOrgResponse,
  GetMyInstanceResponse,
  SetUpOrgRequestSchema,
  SetUpOrgResponse,
} from '@zitadel/proto/zitadel/admin_pb';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { NewAuthService } from './new-auth.service';

@Injectable({
  providedIn: 'root',
})
export class NewAdminService {
  constructor(
    private readonly grpcService: GrpcService,
    private readonly authService: NewAuthService,
  ) {}

  public setupOrg(req: MessageInitShape<typeof SetUpOrgRequestSchema>): Promise<SetUpOrgResponse> {
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
      queryKey: ['admin', 'getMyInstance'],
      queryFn: async () => this.getMyInstance(),
      enabled: (listMyZitadelPermissionsQuery.data() ?? []).includes('iam.write'),
    }));
  }
}
