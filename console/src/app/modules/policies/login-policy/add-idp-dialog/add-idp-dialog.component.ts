import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { IDP, IDPOwnerType, IDPOwnerTypeQuery } from 'src/app/proto/generated/zitadel/idp_pb';
import { IDPQuery } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../../policy-component-types.enum';

@Component({
  selector: 'app-add-idp-dialog',
  templateUrl: './add-idp-dialog.component.html',
  styleUrls: ['./add-idp-dialog.component.scss'],
})
export class AddIdpDialogComponent {
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public idpType!: IDPOwnerType;
  public idpTypes: IDPOwnerType[] = [
    IDPOwnerType.IDP_OWNER_TYPE_SYSTEM,
    IDPOwnerType.IDP_OWNER_TYPE_ORG,
  ];

  public idp: IDP.AsObject | undefined = undefined;
  public availableIdps: Array<IDP.AsObject[] | IDP.AsObject> | string[] = [];

  constructor(
    private mgmtService: ManagementService,
    private toast: ToastService,
    public dialogRef: MatDialogRef<AddIdpDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    if (data.serviceType) {
      this.serviceType = data.serviceType;
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.idpType = IDPOwnerType.IDP_OWNER_TYPE_ORG;
          break;
        case PolicyComponentServiceType.ADMIN:
          this.idpType = IDPOwnerType.IDP_OWNER_TYPE_SYSTEM;
          break;
      }
    }

    this.loadIdps();
  }

  public loadIdps(): void {
    this.idp = undefined;



    switch (this.idpType) {
      case IDPOwnerType.IDP_OWNER_TYPE_ORG:
        const query: IDPQuery = new IDPQuery();
        const idpOTQ: IDPOwnerTypeQuery = new IDPOwnerTypeQuery();
        idpOTQ.setOwnerType(this.idpType);
        query.setOwnerTypeQuery(idpOTQ);

        const orgQuery = new IDPQuery();
        const orgQ = new IDPOwnerTypeQuery();
        orgQ.setOwnerType(IDPOwnerType.IDP_OWNER_TYPE_ORG);
        orgQuery.setOwnerTypeQuery(orgQ);
        this.mgmtService.listOrgIDPs(undefined, undefined, [orgQuery]).then(resp => {
          this.availableIdps = resp.resultList;
        }).catch(error => {
          this.availableIdps = [];
          this.toast.showError(error);
        });
        break;
      case IDPOwnerType.IDP_OWNER_TYPE_SYSTEM:
        const iamQuery = new IDPQuery();
        const iamQ = new IDPOwnerTypeQuery();
        iamQ.setOwnerType(IDPOwnerType.IDP_OWNER_TYPE_SYSTEM);
        iamQuery.setOwnerTypeQuery(iamQ);

        this.mgmtService.listOrgIDPs(undefined, undefined, [iamQuery]).then(resp => {
          this.availableIdps = resp.resultList;
        }).catch(error => {
          this.availableIdps = [];
          this.toast.showError(error);
        });
        break;

    }
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({
      idp: this.idp,
      type: this.idpType,
    });
  }
}
