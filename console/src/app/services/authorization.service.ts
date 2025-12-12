import { GrpcService } from './grpc.service';
import { inject, Injectable } from '@angular/core';
import { NewOrganizationService } from './new-organization.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import { CreateAuthorizationRequestSchema } from '@zitadel/proto/zitadel/authorization/v2beta/authorization_service_pb';
import { mutationOptions } from '@tanstack/angular-query-experimental';

type CreateAuthorizationRequest = Omit<
  Exclude<MessageInitShape<typeof CreateAuthorizationRequestSchema>, { ['$typeName']: string }>,
  'organizationId'
>;

@Injectable({
  providedIn: 'root',
})
export class AuthorizationService {
  private readonly grpcService = inject(GrpcService);
  private readonly orgId = inject(NewOrganizationService).getOrgId();

  private createAuthorization(request: CreateAuthorizationRequest) {
    return this.grpcService.authorization.createAuthorization({ ...request, organizationId: this.orgId() });
  }

  public createAuthorizationMutationOptions = () =>
    mutationOptions({
      mutationKey: ['authorization', 'create'],
      mutationFn: (req: CreateAuthorizationRequest) => this.createAuthorization(req),
    });
}
