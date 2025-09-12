import { Component, Input, OnInit } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { SetDefaultLanguageResponse, SetSecurityPolicyRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { InfoSectionType } from '../../info-section/info-section.component';

@Component({
  selector: 'cnsl-security-policy',
  templateUrl: './security-policy.component.html',
  styleUrls: ['./security-policy.component.scss'],
})
export class SecurityPolicyComponent implements OnInit {
  public originsList: string[] = [];
  public iframeEnabled: boolean = false;
  public impersonationEnabled: boolean = false;

  public loading: boolean = false;
  public InfoSectionType: any = InfoSectionType;

  @Input() public originsControl: UntypedFormControl = new UntypedFormControl({ value: [], disabled: true });

  constructor(
    private service: AdminService,
    private toast: ToastService,
  ) {}

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.service.getSecurityPolicy().then((securityPolicy) => {
      if (securityPolicy.policy) {
        this.impersonationEnabled = securityPolicy.policy?.enableImpersonation;
        this.iframeEnabled = securityPolicy.policy?.enableIframeEmbedding;
        this.originsList = securityPolicy.policy?.allowedOriginsList;
        if (securityPolicy.policy.enableIframeEmbedding) {
          this.originsControl.enable();
        } else {
          this.originsControl.disable();
        }
      }
    });
  }

  private updateData(): Promise<SetDefaultLanguageResponse.AsObject> {
    const req = new SetSecurityPolicyRequest();
    req.setAllowedOriginsList(this.originsList);
    req.setEnableIframeEmbedding(this.iframeEnabled);
    req.setEnableImpersonation(this.impersonationEnabled);
    return (this.service as AdminService).setSecurityPolicy(req);
  }

  public savePolicy(): void {
    const prom = this.updateData();
    this.loading = true;
    if (prom) {
      prom
        .then(() => {
          this.toast.showInfo('POLICY.SECURITY_POLICY.SAVED', true);
          this.loading = false;
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.loading = false;
          this.toast.showError(error);
        });
    }
  }

  public add(input: any): void {
    if (this.originsControl.valid) {
      if (input.value !== '' && input.value !== ' ' && input.value !== '/') {
        this.originsList.push(input.value);
      }
      if (input) {
        input.value = '';
      }
    }
  }

  public remove(redirect: any): void {
    const index = this.originsList.indexOf(redirect);

    if (index >= 0) {
      this.originsList.splice(index, 1);
    }
  }

  public iframeEnabledChanged(event: MatCheckboxChange) {
    if (event.checked) {
      this.originsControl.enable();
    } else {
      this.originsControl.disable();
    }
  }
}
