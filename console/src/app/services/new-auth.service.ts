import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { create } from '@bufbuild/protobuf';
import {
  AddMyAuthFactorOTPSMSResponse,
  GetMyLoginPolicyResponse,
  GetMyLoginPolicyRequestSchema,
  GetMyPasswordComplexityPolicyResponse,
  GetMyUserResponse,
  ListMyAuthFactorsRequestSchema,
  ListMyAuthFactorsResponse,
  RemoveMyAuthFactorOTPEmailRequestSchema,
  RemoveMyAuthFactorOTPEmailResponse,
  RemoveMyAuthFactorOTPRequestSchema,
  RemoveMyAuthFactorOTPResponse,
  RemoveMyAuthFactorU2FRequestSchema,
  RemoveMyAuthFactorU2FResponse,
  RemoveMyAuthFactorOTPSMSRequestSchema,
  RemoveMyAuthFactorOTPSMSResponse,
  ListMyMetadataResponse,
  VerifyMyPhoneResponse,
} from '@zitadel/proto/zitadel/auth_pb';

@Injectable({
  providedIn: 'root',
})
export class NewAuthService {
  constructor(private readonly grpcService: GrpcService) {}

  public getMyUser(): Promise<GetMyUserResponse> {
    return this.grpcService.authNew.getMyUser({});
  }

  public verifyMyPhone(code: string): Promise<VerifyMyPhoneResponse> {
    return this.grpcService.authNew.verifyMyPhone({ code });
  }

  public addMyAuthFactorOTPSMS(): Promise<AddMyAuthFactorOTPSMSResponse> {
    return this.grpcService.authNew.addMyAuthFactorOTPSMS({});
  }

  public listMyMetadata(): Promise<ListMyMetadataResponse> {
    return this.grpcService.authNew.listMyMetadata({});
  }

  public listMyMultiFactors(): Promise<ListMyAuthFactorsResponse> {
    return this.grpcService.authNew.listMyAuthFactors(create(ListMyAuthFactorsRequestSchema), null);
  }

  public removeMyAuthFactorOTPSMS(): Promise<RemoveMyAuthFactorOTPSMSResponse> {
    return this.grpcService.authNew.removeMyAuthFactorOTPSMS(create(RemoveMyAuthFactorOTPSMSRequestSchema), null);
  }

  public getMyLoginPolicy(): Promise<GetMyLoginPolicyResponse> {
    return this.grpcService.authNew.getMyLoginPolicy(create(GetMyLoginPolicyRequestSchema), null);
  }

  public removeMyMultiFactorOTP(): Promise<RemoveMyAuthFactorOTPResponse> {
    return this.grpcService.authNew.removeMyAuthFactorOTP(create(RemoveMyAuthFactorOTPRequestSchema), null);
  }

  public removeMyMultiFactorU2F(tokenId: string): Promise<RemoveMyAuthFactorU2FResponse> {
    return this.grpcService.authNew.removeMyAuthFactorU2F(create(RemoveMyAuthFactorU2FRequestSchema, { tokenId }), null);
  }

  public removeMyAuthFactorOTPEmail(): Promise<RemoveMyAuthFactorOTPEmailResponse> {
    return this.grpcService.authNew.removeMyAuthFactorOTPEmail(create(RemoveMyAuthFactorOTPEmailRequestSchema), null);
  }

  public getMyPasswordComplexityPolicy(): Promise<GetMyPasswordComplexityPolicyResponse> {
    return this.grpcService.authNew.getMyPasswordComplexityPolicy({});
  }
}
