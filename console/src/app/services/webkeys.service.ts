import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
// todo: this should be from v2 but local generation is a bit funky
import {
  ActivateWebKeyResponse,
  CreateWebKeyRequestSchema,
  CreateWebKeyResponse,
  DeleteWebKeyResponse,
  ListWebKeysResponse,
} from '@zitadel/proto/zitadel/resources/webkey/v3alpha/webkey_service_pb';
import type { MessageInitShape } from '@bufbuild/protobuf';

@Injectable({
  providedIn: 'root',
})
export class WebKeysService {
  constructor(private readonly grpcService: GrpcService) {}

  public ListWebKeys(): Promise<ListWebKeysResponse> {
    return this.grpcService.webKeysNew.listWebKeys({});
  }

  public DeleteWebKey(id: string): Promise<DeleteWebKeyResponse> {
    return this.grpcService.webKeysNew.deleteWebKey({ id });
  }

  public CreateWebKey(req: MessageInitShape<typeof CreateWebKeyRequestSchema>): Promise<CreateWebKeyResponse> {
    return this.grpcService.webKeysNew.createWebKey(req);
  }

  public ActivateWebKey(id: string): Promise<ActivateWebKeyResponse> {
    return this.grpcService.webKeysNew.activateWebKey({ id });
  }
}
