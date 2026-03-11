import { Injectable } from '@angular/core';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  GetDetectionSettingsResponse,
  SetDetectionSettingsRequestSchema,
  SetDetectionSettingsResponse,
} from '@zitadel/proto/zitadel/settings/v2/detection_settings_pb';
import {
  CreateDetectionRuleRequestSchema,
  CreateDetectionRuleResponse,
  DeleteDetectionRuleResponse,
  GetDetectionRuleResponse,
  ListDetectionRulesResponse,
  UpdateDetectionRuleRequestSchema,
  UpdateDetectionRuleResponse,
} from '@zitadel/proto/zitadel/settings/v2/detection_rules_pb';

import { GrpcService } from './grpc.service';

@Injectable({
  providedIn: 'root',
})
export class NewSettingsService {
  constructor(private readonly grpcService: GrpcService) {}

  public getDetectionSettings(): Promise<GetDetectionSettingsResponse> {
    return this.grpcService.settingsNew.getDetectionSettings({});
  }

  public setDetectionSettings(
    req: MessageInitShape<typeof SetDetectionSettingsRequestSchema>,
  ): Promise<SetDetectionSettingsResponse> {
    return this.grpcService.settingsNew.setDetectionSettings(req);
  }

  public listDetectionRules(): Promise<ListDetectionRulesResponse> {
    return this.grpcService.settingsNew.listDetectionRules({});
  }

  public getDetectionRule(ruleId: string): Promise<GetDetectionRuleResponse> {
    return this.grpcService.settingsNew.getDetectionRule({ ruleId });
  }

  public createDetectionRule(
    req: MessageInitShape<typeof CreateDetectionRuleRequestSchema>,
  ): Promise<CreateDetectionRuleResponse> {
    return this.grpcService.settingsNew.createDetectionRule(req);
  }

  public updateDetectionRule(
    req: MessageInitShape<typeof UpdateDetectionRuleRequestSchema>,
  ): Promise<UpdateDetectionRuleResponse> {
    return this.grpcService.settingsNew.updateDetectionRule(req);
  }

  public deleteDetectionRule(ruleId: string): Promise<DeleteDetectionRuleResponse> {
    return this.grpcService.settingsNew.deleteDetectionRule({ ruleId });
  }
}
