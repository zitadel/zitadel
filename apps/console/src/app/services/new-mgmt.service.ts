import { Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import {
  GenerateMachineSecretRequestSchema,
  GenerateMachineSecretResponse,
  GetDefaultPasswordComplexityPolicyResponse,
  GetLoginPolicyRequestSchema,
  GetLoginPolicyResponse,
  GetPasswordComplexityPolicyResponse,
  ListUserMetadataRequestSchema,
  ListUserMetadataResponse,
  RemoveMachineSecretRequestSchema,
  RemoveMachineSecretResponse,
  RemoveUserMetadataRequestSchema,
  RemoveUserMetadataResponse,
  ResendHumanEmailVerificationRequestSchema,
  ResendHumanEmailVerificationResponse,
  ResendHumanInitializationRequestSchema,
  ResendHumanInitializationResponse,
  ResendHumanPhoneVerificationRequestSchema,
  ResendHumanPhoneVerificationResponse,
  SendHumanResetPasswordNotificationRequest_Type,
  SendHumanResetPasswordNotificationRequestSchema,
  SendHumanResetPasswordNotificationResponse,
  SetUserMetadataRequestSchema,
  SetUserMetadataResponse,
  UpdateMachineRequestSchema,
  UpdateMachineResponse,
} from '@zitadel/proto/zitadel/management_pb';
import { MessageInitShape, create } from '@bufbuild/protobuf';

@Injectable({
  providedIn: 'root',
})
export class NewMgmtService {
  constructor(private readonly grpcService: GrpcService) {}

  public getLoginPolicy(): Promise<GetLoginPolicyResponse> {
    return this.grpcService.mgmtNew.getLoginPolicy(create(GetLoginPolicyRequestSchema));
  }

  public generateMachineSecret(userId: string): Promise<GenerateMachineSecretResponse> {
    return this.grpcService.mgmtNew.generateMachineSecret(create(GenerateMachineSecretRequestSchema, { userId }));
  }

  public removeMachineSecret(userId: string): Promise<RemoveMachineSecretResponse> {
    return this.grpcService.mgmtNew.removeMachineSecret(create(RemoveMachineSecretRequestSchema, { userId }));
  }

  public updateMachine(req: MessageInitShape<typeof UpdateMachineRequestSchema>): Promise<UpdateMachineResponse> {
    return this.grpcService.mgmtNew.updateMachine(create(UpdateMachineRequestSchema, req));
  }

  public resendHumanEmailVerification(userId: string): Promise<ResendHumanEmailVerificationResponse> {
    return this.grpcService.mgmtNew.resendHumanEmailVerification(
      create(ResendHumanEmailVerificationRequestSchema, { userId }),
    );
  }

  public resendHumanPhoneVerification(userId: string): Promise<ResendHumanPhoneVerificationResponse> {
    return this.grpcService.mgmtNew.resendHumanPhoneVerification(
      create(ResendHumanPhoneVerificationRequestSchema, { userId }),
    );
  }

  public sendHumanResetPasswordNotification(
    req: MessageInitShape<typeof SendHumanResetPasswordNotificationRequestSchema>,
  ): Promise<SendHumanResetPasswordNotificationResponse> {
    return this.grpcService.mgmtNew.sendHumanResetPasswordNotification(
      create(SendHumanResetPasswordNotificationRequestSchema, req),
    );
  }

  public resendHumanInitialization(userId: string, email: string = ''): Promise<ResendHumanInitializationResponse> {
    return this.grpcService.mgmtNew.resendHumanInitialization(
      create(ResendHumanInitializationRequestSchema, { userId, email }),
    );
  }

  public listUserMetadata(id: string): Promise<ListUserMetadataResponse> {
    return this.grpcService.mgmtNew.listUserMetadata(create(ListUserMetadataRequestSchema, { id }));
  }

  public setUserMetadata(req: MessageInitShape<typeof SetUserMetadataRequestSchema>): Promise<SetUserMetadataResponse> {
    return this.grpcService.mgmtNew.setUserMetadata(create(SetUserMetadataRequestSchema, req));
  }

  public removeUserMetadata(
    req: MessageInitShape<typeof RemoveUserMetadataRequestSchema>,
  ): Promise<RemoveUserMetadataResponse> {
    return this.grpcService.mgmtNew.removeUserMetadata(create(RemoveUserMetadataRequestSchema, req));
  }

  public getPasswordComplexityPolicy(): Promise<GetPasswordComplexityPolicyResponse> {
    return this.grpcService.mgmtNew.getPasswordComplexityPolicy({});
  }

  public getDefaultPasswordComplexityPolicy(): Promise<GetDefaultPasswordComplexityPolicyResponse> {
    return this.grpcService.mgmtNew.getDefaultPasswordComplexityPolicy({});
  }
}
