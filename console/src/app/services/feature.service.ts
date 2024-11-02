import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';

import {
  GetInstanceFeaturesRequest,
  GetInstanceFeaturesResponse,
  ResetInstanceFeaturesRequest,
  SetInstanceFeaturesRequest,
  SetInstanceFeaturesResponse,
} from '../proto/generated/zitadel/feature/v2beta/instance_pb';
import {
  GetOrganizationFeaturesRequest,
  GetOrganizationFeaturesResponse,
} from '../proto/generated/zitadel/feature/v2beta/organization_pb';
import { GetUserFeaturesRequest, GetUserFeaturesResponse } from '../proto/generated/zitadel/feature/v2beta/user_pb';
import { GetSystemFeaturesRequest, GetSystemFeaturesResponse } from '../proto/generated/zitadel/feature/v2beta/system_pb';

@Injectable({
  providedIn: 'root',
})
export class FeatureService {
  constructor(private readonly grpcService: GrpcService) {}

  public getInstanceFeatures(inheritance: boolean): Promise<GetInstanceFeaturesResponse> {
    const req = new GetInstanceFeaturesRequest();
    req.setInheritance(inheritance);
    return this.grpcService.feature.getInstanceFeatures(req, null).then((resp) => resp);
  }

  public setInstanceFeatures(req: SetInstanceFeaturesRequest): Promise<SetInstanceFeaturesResponse> {
    return this.grpcService.feature.setInstanceFeatures(req, null);
  }

  public resetInstanceFeatures(): Promise<SetInstanceFeaturesResponse> {
    const req = new ResetInstanceFeaturesRequest();
    return this.grpcService.feature.resetInstanceFeatures(req, null);
  }

  public getOrganizationFeatures(orgId: string, inheritance: boolean): Promise<GetOrganizationFeaturesResponse> {
    const req = new GetOrganizationFeaturesRequest();
    req.setOrganizationId(orgId);
    req.setInheritance(inheritance);
    return this.grpcService.feature.getOrganizationFeatures(req, null).then((resp) => resp);
  }

  public getSystemFeatures(): Promise<GetSystemFeaturesResponse> {
    const req = new GetSystemFeaturesRequest();
    return this.grpcService.feature.getSystemFeatures(req, null).then((resp) => resp);
  }

  public getUserFeatures(userId: string, inheritance: boolean): Promise<GetUserFeaturesResponse> {
    const req = new GetUserFeaturesRequest();
    req.setInheritance(inheritance);
    req.setUserId(userId);
    return this.grpcService.feature.getUserFeatures(req, null).then((resp) => resp);
  }
}
