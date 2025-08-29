import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import type { MessageInitShape } from '@bufbuild/protobuf';
import {
  DeleteWebKeyResponse,
  ListWebKeysResponse,
  CreateWebKeyRequestSchema,
  CreateWebKeyResponse,
  ActivateWebKeyResponse,
} from '@zitadel/proto/zitadel/webkey/v2beta/webkey_service_pb';

@Injectable({
  providedIn: 'root',
})
export class WebKeysService {
  constructor(private readonly grpcService: GrpcService) {}

  public ListWebKeys(): Promise<ListWebKeysResponse> {
    return this.grpcService.webKey['listWebKeys']({});
  }

  public DeleteWebKey(id: string): Promise<DeleteWebKeyResponse> {
    return this.grpcService.webKey['deleteWebKey']({ id });
  }

  public CreateWebKey(req: MessageInitShape<typeof CreateWebKeyRequestSchema>): Promise<CreateWebKeyResponse> {
    return this.grpcService.webKey['createWebKey'](req);
  }

  public ActivateWebKey(id: string): Promise<ActivateWebKeyResponse> {
    return this.grpcService.webKey['activateWebKey']({ id });
  }
}
