import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { create } from '@bufbuild/protobuf';
import {
  AddMyAuthFactorOTPSMSRequestSchema,
  AddMyAuthFactorOTPSMSResponse,
  GetMyUserRequestSchema,
  GetMyUserResponse,
  VerifyMyPhoneRequestSchema,
  VerifyMyPhoneResponse,
} from '@zitadel/proto/zitadel/auth_pb';

@Injectable({
  providedIn: 'root',
})
export class NewAuthService {
  constructor(private readonly grpcService: GrpcService) {}

  public getMyUser(): Promise<GetMyUserResponse> {
    return this.grpcService.authNew.getMyUser(create(GetMyUserRequestSchema));
  }

  public verifyMyPhone(code: string): Promise<VerifyMyPhoneResponse> {
    return this.grpcService.authNew.verifyMyPhone(create(VerifyMyPhoneRequestSchema, { code }));
  }

  public addMyAuthFactorOTPSMS(): Promise<AddMyAuthFactorOTPSMSResponse> {
    return this.grpcService.authNew.addMyAuthFactorOTPSMS(create(AddMyAuthFactorOTPSMSRequestSchema));
  }
}
