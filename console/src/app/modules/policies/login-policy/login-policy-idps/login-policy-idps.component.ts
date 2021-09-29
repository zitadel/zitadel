import { Component, Input, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { IDP, IDPLoginPolicyLink, IDPOwnerType, IDPStylingType } from 'src/app/proto/generated/zitadel/idp_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../../policy-component-types.enum';
import { AddIdpDialogComponent } from './add-idp-dialog/add-idp-dialog.component';

@Component({
  selector: 'cnsl-login-policy-idps',
  templateUrl: './login-policy-idps.component.html',
  styleUrls: ['./login-policy-idps.component.scss']
})
export class LoginPolicyIdpsComponent implements OnInit {
  @Input() public disabled: boolean = true;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  @Input() public service!: ManagementService | AdminService;
  public loading: boolean = false;

  public idps: IDPLoginPolicyLink.AsObject[] = [];

  public IDPStylingType: any = IDPStylingType;

  constructor(
    private toast: ToastService,
    private dialog: MatDialog,
  ) { }

  ngOnInit(): void {
    this.getIdps().then(resp => {
      this.idps = resp;
      console.log(this.idps);
    });
  }

  private async getIdps(): Promise<IDPLoginPolicyLink.AsObject[]> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).listLoginPolicyIDPs()
          .then((resp) => {
            return resp.resultList;
          });
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).listLoginPolicyIDPs()
          .then((providers) => {
            return providers.resultList;
          });
    }
  }

  private addIdp(idp: IDP.AsObject | IDP.AsObject, ownerType: IDPOwnerType): Promise<any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).addIDPToLoginPolicy(idp.id, ownerType);
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).addIDPToLoginPolicy(idp.id);
    }
  }

  public openDialog(): void {
    const dialogRef = this.dialog.open(AddIdpDialogComponent, {
      data: {
        serviceType: this.serviceType,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp && resp.idp && resp.type) {
        this.addIdp(resp.idp, resp.type).then(() => {
          this.loading = true;
          setTimeout(() => {
            this.getIdps();
          }, 1000);
        }).catch(error => {
          this.toast.showError(error);
        });
      }
    });
  }

  public removeIdp(idp: IDPLoginPolicyLink.AsObject): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        (this.service as ManagementService).removeIDPFromLoginPolicy(idp.idpId).then(() => {
          const index = this.idps.findIndex(temp => temp === idp);
          if (index > -1) {
            this.idps.splice(index, 1);
          }
        }, error => {
          this.toast.showError(error);
        });
        break;
      case PolicyComponentServiceType.ADMIN:
        (this.service as AdminService).removeIDPFromLoginPolicy(idp.idpId).then(() => {
          const index = this.idps.findIndex(temp => temp === idp);
          if (index > -1) {
            this.idps.splice(index, 1);
          }
        }, error => {
          this.toast.showError(error);
        });
        break;
    }
  }
}
