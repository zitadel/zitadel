import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { SMTPKnownProviders } from '../../smtp-provider/known-smtp-providers-settings';

@Component({
  selector: 'cnsl-smtp-settings',
  templateUrl: './smtp-settings.component.html',
  styleUrls: ['./smtp-settings.component.scss'],
})
export class SMTPSettingsComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  public service!: ManagementService | AdminService;

  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public providers = SMTPKnownProviders;

  constructor(private injector: Injector) {}

  ngOnInit(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        this.service = this.injector.get(ManagementService as Type<ManagementService>);
        break;
      case PolicyComponentServiceType.ADMIN:
        this.service = this.injector.get(AdminService as Type<AdminService>);
        break;
    }
  }

  public createGoogle() {}
}
