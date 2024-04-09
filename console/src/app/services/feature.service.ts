import { Injectable } from '@angular/core';
import { BehaviorSubject, catchError, finalize, from, map, Observable, of, Subject, switchMap, tap } from 'rxjs';

import { FeatureFlag } from '../proto/generated/zitadel/feature/v2beta/feature_pb';
import { Event } from '../proto/generated/zitadel/event_pb';
import {
  ResetCustomDomainClaimedMessageTextToDefaultRequest,
  ResetCustomDomainClaimedMessageTextToDefaultResponse,
  ResetCustomInitMessageTextToDefaultRequest,
  ResetCustomInitMessageTextToDefaultResponse,
  ResetCustomPasswordChangeMessageTextToDefaultRequest,
  ResetCustomPasswordChangeMessageTextToDefaultResponse,
  ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest,
  ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse,
  ResetCustomPasswordResetMessageTextToDefaultRequest,
  ResetCustomPasswordResetMessageTextToDefaultResponse,
  ResetCustomVerifyEmailMessageTextToDefaultRequest,
  ResetCustomVerifyEmailMessageTextToDefaultResponse,
  ResetCustomVerifyEmailOTPMessageTextToDefaultRequest,
  ResetCustomVerifyEmailOTPMessageTextToDefaultResponse,
  ResetCustomVerifyPhoneMessageTextToDefaultRequest,
  ResetCustomVerifyPhoneMessageTextToDefaultResponse,
  ResetCustomVerifySMSOTPMessageTextToDefaultRequest,
  ResetCustomVerifySMSOTPMessageTextToDefaultResponse,
} from '../proto/generated/zitadel/management_pb';
import { SearchQuery } from '../proto/generated/zitadel/member_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { GrpcService } from './grpc.service';
import { StorageLocation, StorageService } from './storage.service';
import {
  IsReachedQuery,
  Milestone,
  MilestoneQuery,
  MilestoneType,
} from '../proto/generated/zitadel/milestone/v1/milestone_pb';
import { OrgFieldName, OrgQuery } from '../proto/generated/zitadel/org_pb';
import { SortDirection } from '@angular/material/sort';
import {
  GetInstanceFeaturesRequest,
  GetInstanceFeaturesResponse,
} from '../proto/generated/zitadel/feature/v2beta/instance_pb';
import {
  GetOrganizationFeaturesRequest,
  GetOrganizationFeaturesResponse,
} from '../proto/generated/zitadel/feature/v2beta/organization_pb';
import { GetUserFeaturesRequest, GetUserFeaturesResponse } from '../proto/generated/zitadel/feature/v2beta/user_pb';
import { GetSystemFeaturesRequest, GetSystemFeaturesResponse } from '../proto/generated/zitadel/feature/v2beta/system_pb';
import { inherits } from 'util';

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
