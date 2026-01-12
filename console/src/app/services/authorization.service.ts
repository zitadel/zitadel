import { GrpcService } from './grpc.service';
import { inject, Injectable } from '@angular/core';
import { MessageInitShape } from '@bufbuild/protobuf';
import { CreateAuthorizationRequestSchema } from '@zitadel/proto/zitadel/authorization/v2beta/authorization_service_pb';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { StorageKey, StorageLocation, StorageService } from './storage.service';

type CreateAuthorizationRequest = Omit<
  Exclude<MessageInitShape<typeof CreateAuthorizationRequestSchema>, { ['$typeName']: string }>,
  'organizationId'
>;

@Injectable({
  providedIn: 'root',
})
export class AuthorizationService {
  private readonly grpcService = inject(GrpcService);
  private readonly storageService = inject(StorageService);

  private createAuthorization(request: CreateAuthorizationRequest) {
    const organizationId = this.storageService.getItem(StorageKey.organizationId, StorageLocation.session) ?? undefined;
    return this.grpcService.authorization.createAuthorization({ ...request, organizationId });
  }

  public createAuthorizationMutationOptions = () =>
    mutationOptions({
      mutationKey: ['authorization', 'create'],
      mutationFn: (req: CreateAuthorizationRequest) => this.createAuthorization(req),
    });
}
