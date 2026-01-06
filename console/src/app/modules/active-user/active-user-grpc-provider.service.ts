import { GrpcService } from "@/services/grpc.service";
import { Client } from "@connectrpc/connect";
import { ActiveUserService as ActiveUserServiceGrpc } from "@zitadel/proto/zitadel/analytics/v2beta/active_user_service_pb";
import { inject, Injectable } from "@angular/core";
import { ActiveUserGrpcMockService } from "./active-user-grpc-mock.service";

@Injectable()
export class ActiveUserGrpcProviderService {
  private readonly grpcService = inject(GrpcService);

  public getClient(): Client<typeof ActiveUserServiceGrpc> {
    return this.grpcService.activeUser;
  }
}

@Injectable()
export class ActiveUserGrpcMockProviderService extends ActiveUserGrpcProviderService {
  public override getClient(): Client<typeof ActiveUserServiceGrpc> {
    return new ActiveUserGrpcMockService();
  }
}
