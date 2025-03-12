import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  GetInstanceFeaturesResponse,
  SetInstanceFeaturesRequest,
  SetInstanceFeaturesResponse,
} from '@zitadel/proto/zitadel/feature/v2/instance_pb';
import { ResetInstanceFeaturesRequest } from '../proto/generated/zitadel/feature/v2/instance_pb';

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

  public resetInstanceFeatures(): Promise<SetInstanceFeaturesResponse> {
    const req = new ResetInstanceFeaturesRequest();
    return this.grpcService.featureNew.resetInstanceFeatures(req);
  }
}
