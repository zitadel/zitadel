import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { IDP, IDPOwnerType } from 'src/app/proto/generated/zitadel/idp_pb';
import { IDPQuery } from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';

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
    public IdpProviderType: any = IDPOwnerType;

    constructor(
        private mgmtService: ManagementService,
        private adminService: AdminService,
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
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            const query: IDPQuery = new IDPQuery();
            query.set(IdpSearchKey.IDPSEARCHKEY_PROVIDER_TYPE);
            query.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
            query.setValue(this.idpType.toString());

            this.mgmtService.listOrgIDPs(undefined, undefined, [query]).then(resp => {
                this.availableIdps = resp.resultList;
            });
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
            this.adminService.listIDPs().then(resp => {
                this.availableIdps = resp.resultList;
            });
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
