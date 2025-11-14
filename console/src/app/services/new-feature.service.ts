import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  GetInstanceFeaturesResponse,
  ResetInstanceFeaturesResponse,
  SetInstanceFeaturesRequestSchema,
  SetInstanceFeaturesResponse,
} from '@zitadel/proto/zitadel/feature/v2/instance_pb';
import { MessageInitShape } from '@bufbuild/protobuf';

@Injectable({
  providedIn: 'root',
})
export class NewFeatureService {
  constructor(private readonly grpcService: GrpcService) {}

  public getInstanceFeatures(): Promise<GetInstanceFeaturesResponse> {
    return this.grpcService.featureNew.getInstanceFeatures({});
  }

  public setInstanceFeatures(
    req: MessageInitShape<typeof SetInstanceFeaturesRequestSchema>,
  ): Promise<SetInstanceFeaturesResponse> {
    return this.grpcService.featureNew.setInstanceFeatures(req);
  }

  public resetInstanceFeatures(): Promise<ResetInstanceFeaturesResponse> {
    return this.grpcService.featureNew.resetInstanceFeatures({});
  }
}
