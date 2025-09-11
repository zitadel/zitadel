import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  CreateTargetRequestSchema,
  CreateTargetResponse,
  DeleteTargetRequestSchema,
  DeleteTargetResponse,
  GetTargetRequestSchema,
  GetTargetResponse,
  ListExecutionFunctionsRequestSchema,
  ListExecutionFunctionsResponse,
  ListExecutionMethodsRequestSchema,
  ListExecutionMethodsResponse,
  ListExecutionServicesRequestSchema,
  ListExecutionServicesResponse,
  ListExecutionsRequestSchema,
  ListExecutionsResponse,
  ListTargetsRequestSchema,
  ListTargetsResponse,
  SetExecutionRequestSchema,
  SetExecutionResponse,
  UpdateTargetRequestSchema,
  UpdateTargetResponse,
} from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';
import { UserService } from './user.service';
import { injectQuery } from '@tanstack/angular-query-experimental';

@Injectable({
  providedIn: 'root',
})
export class ActionService {
  constructor(
    private readonly grpcService: GrpcService,
    private userService: UserService,
  ) {}

  public listTargets(req: MessageInitShape<typeof ListTargetsRequestSchema>): Promise<ListTargetsResponse> {
    return this.grpcService.actionNew.listTargets(req);
  }

  public createTarget(req: MessageInitShape<typeof CreateTargetRequestSchema>): Promise<CreateTargetResponse> {
    return this.grpcService.actionNew.createTarget(req);
  }

  public deleteTarget(req: MessageInitShape<typeof DeleteTargetRequestSchema>): Promise<DeleteTargetResponse> {
    return this.grpcService.actionNew.deleteTarget(req);
  }

  public getTarget(req: MessageInitShape<typeof GetTargetRequestSchema>): Promise<GetTargetResponse> {
    return this.grpcService.actionNew.getTarget(req);
  }

  public updateTarget(req: MessageInitShape<typeof UpdateTargetRequestSchema>): Promise<UpdateTargetResponse> {
    return this.grpcService.actionNew.updateTarget(req);
  }

  public listExecutions(req: MessageInitShape<typeof ListExecutionsRequestSchema>): Promise<ListExecutionsResponse> {
    return this.grpcService.actionNew.listExecutions(req);
  }

  public setExecution(req: MessageInitShape<typeof SetExecutionRequestSchema>): Promise<SetExecutionResponse> {
    return this.grpcService.actionNew.setExecution(req);
  }

  private listExecutionServices(
    req: MessageInitShape<typeof ListExecutionServicesRequestSchema>,
  ): Promise<ListExecutionServicesResponse> {
    return this.grpcService.actionNew.listExecutionServices(req);
  }

  public listExecutionServicesQuery() {
    return injectQuery(() => ({
      queryKey: [this.userService.userId(), 'action', 'listExecutionServices'],
      queryFn: () => this.listExecutionServices({}).then(({ services }) => services),
    }));
  }

  private listExecutionMethods(
    req: MessageInitShape<typeof ListExecutionMethodsRequestSchema>,
  ): Promise<ListExecutionMethodsResponse> {
    return this.grpcService.actionNew.listExecutionMethods(req);
  }

  public listExecutionMethodsQuery() {
    return injectQuery(() => ({
      queryKey: [this.userService.userId(), 'action', 'listExecutionMethods'],
      queryFn: () => this.listExecutionMethods({}).then(({ methods }) => methods),
    }));
  }

  private listExecutionFunctions(
    req: MessageInitShape<typeof ListExecutionFunctionsRequestSchema>,
  ): Promise<ListExecutionFunctionsResponse> {
    return this.grpcService.actionNew.listExecutionFunctions(req);
  }

  public listExecutionFunctionsQuery() {
    return injectQuery(() => ({
      queryKey: [this.userService.userId(), 'action', 'listExecutionFunctions'],
      queryFn: () => this.listExecutionFunctions({}).then(({ functions }) => functions),
    }));
  }
}
