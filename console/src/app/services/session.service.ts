import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import type { MessageInitShape } from '@bufbuild/protobuf';
import { ListMyUserSessionsResponse, ListMyUserSessionsResponseSchema } from '@zitadel/proto/zitadel/auth_pb';

@Injectable({
  providedIn: 'root',
})
export class SessionService {
  constructor(private readonly grpcService: GrpcService) {}

  public listMyUserSessions(
    req: MessageInitShape<typeof ListMyUserSessionsResponseSchema>,
  ): Promise<ListMyUserSessionsResponse> {
    return this.grpcService.session.listMyUserSessions(req);
  }
}
