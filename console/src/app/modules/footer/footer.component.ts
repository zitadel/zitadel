import { Component, Input } from '@angular/core';
import { LabelPolicy, PrivacyPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

@Component({
  selector: 'cnsl-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss'],
})
export class FooterComponent {
  public policy!: PrivacyPolicy.AsObject;
  @Input() public privateLabelPolicy!: LabelPolicy.AsObject;
  constructor(mgmtService: ManagementService) {
    mgmtService.getPrivacyPolicy().then((policyResp) => {
      if (policyResp.policy) {
        this.policy = policyResp.policy;
      }
    });
  }
}
