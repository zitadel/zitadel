import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  AddMyAuthFactorOTPSMSResponse,
  GetMyPasswordComplexityPolicyResponse,
  GetMyUserResponse,
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

  public getMyPasswordComplexityPolicy(): Promise<GetMyPasswordComplexityPolicyResponse> {
    return this.grpcService.authNew.getMyPasswordComplexityPolicy({});
  }
}
