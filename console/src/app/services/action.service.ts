import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  CreateTargetRequestSchema,
  CreateTargetResponse,
  GetTargetResponse,
  ListExecutionFunctionsResponse,
  ListExecutionServicesResponse,
  SearchExecutionsRequestSchema,
  SearchExecutionsResponse,
  SearchTargetsRequestSchema,
  SearchTargetsResponse,
} from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';

@Injectable({
  providedIn: 'root',
})
export class ActionService {
  constructor(private readonly grpcService: GrpcService) {}

  public searchTargets(req: MessageInitShape<typeof SearchTargetsRequestSchema>): Promise<SearchTargetsResponse> {
    return this.grpcService.actionNew.searchTargets(req);
  }

  public getTarget(id: string): Promise<GetTargetResponse> {
    return this.grpcService.actionNew.getTarget({ id });
  }

  public createTarget(req: MessageInitShape<typeof CreateTargetRequestSchema>): Promise<CreateTargetResponse> {
    return this.grpcService.actionNew.createTarget(req);
  }

  public searchExections(req: MessageInitShape<typeof SearchExecutionsRequestSchema>): Promise<SearchExecutionsResponse> {
    return this.grpcService.actionNew.searchExecutions(req);
  }

  public listExecutionFunctions(): Promise<ListExecutionFunctionsResponse> {
    return this.grpcService.actionNew.listExecutionFunctions({});
  }

  public listExecutionServices(): Promise<ListExecutionServicesResponse> {
    return this.grpcService.actionNew.listExecutionServices({});
  }
}
