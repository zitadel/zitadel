import { Injectable } from '@angular/core';
import { MessageInitShape } from '@bufbuild/protobuf';
import { UpdateApplicationRequestSchema } from '@zitadel/proto/zitadel/application/v2/application_service_pb';

import { GrpcService } from './grpc.service';

@Injectable({
  providedIn: 'root',
})
export class ApplicationService {
  constructor(private readonly grpcService: GrpcService) {}

  public updateApplication(req: MessageInitShape<typeof UpdateApplicationRequestSchema>) {
    return this.grpcService.application.updateApplication(req);
  }
}
