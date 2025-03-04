import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { GetInstanceFeaturesResponse } from '@zitadel/proto/zitadel/feature/v2/instance_pb';

@Injectable({
  providedIn: 'root',
})
export class NewFeatureService {
  constructor(private readonly grpcService: GrpcService) {}

  public getInstanceFeatures(): Promise<GetInstanceFeaturesResponse> {
    return this.grpcService.featureNew.getInstanceFeatures({});
  }
}
