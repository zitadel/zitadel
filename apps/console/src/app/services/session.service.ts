import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import type { MessageInitShape } from '@bufbuild/protobuf';
import { ListSessionsRequestSchema, ListSessionsResponse } from '@zitadel/proto/zitadel/session/v2/session_service_pb';

@Injectable({
  providedIn: 'root',
})
export class SessionService {
  constructor(private readonly grpcService: GrpcService) {}

  public listSessions(req: MessageInitShape<typeof ListSessionsRequestSchema>): Promise<ListSessionsResponse> {
    return this.grpcService.session.listSessions(req);
  }
}
