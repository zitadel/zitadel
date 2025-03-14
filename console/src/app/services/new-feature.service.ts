import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  GetInstanceFeaturesResponse,
  ResetInstanceFeaturesResponse,
  SetInstanceFeaturesRequest,
  SetInstanceFeaturesResponse,
} from '@zitadel/proto/zitadel/feature/v2/instance_pb';

@Injectable({
  providedIn: 'root',
})
export class NewFeatureService {
  constructor(private readonly grpcService: GrpcService) {}

  public getInstanceFeatures(): Promise<GetInstanceFeaturesResponse> {
    return this.grpcService.featureNew.getInstanceFeatures({});
  }

  public setInstanceFeatures(req: SetInstanceFeaturesRequest): Promise<SetInstanceFeaturesResponse> {
    return this.grpcService.featureNew.setInstanceFeatures(req);
  }

  public resetInstanceFeatures(): Promise<ResetInstanceFeaturesResponse> {
    return this.grpcService.featureNew.resetInstanceFeatures({});
  }
}
