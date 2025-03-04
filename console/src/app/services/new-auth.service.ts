import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  AddMyAuthFactorOTPSMSResponse,
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
    return this.grpcService.authNew.verifyMyPhone({});
  }

  public addMyAuthFactorOTPSMS(): Promise<AddMyAuthFactorOTPSMSResponse> {
    return this.grpcService.authNew.addMyAuthFactorOTPSMS({});
  }

  public listMyMetadata(): Promise<ListMyMetadataResponse> {
    return this.grpcService.authNew.listMyMetadata({});
  }
}
