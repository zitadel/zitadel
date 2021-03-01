import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { IDP } from 'src/app/proto/generated/zitadel/idp_pb';
import {
    Idp,
    IdpProviderType,
    IdpSearchKey,
    IdpSearchQuery,
    IdpView as MgmtIdpView,
    SearchMethod,
} from 'src/app/proto/generated/zitadel/management_pb';
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

    public idpType!: IdpProviderType;
    public idpTypes: IdpProviderType[] = [
        IdpProviderType.IDPPROVIDERTYPE_SYSTEM,
        IdpProviderType.IDPPROVIDERTYPE_ORG,
    ];

    public idp: Idp.AsObject | undefined = undefined;
    public availableIdps: Array<IDP.AsObject[] | MgmtIdpView.AsObject> | string[] = [];
    public IdpProviderType: any = IdpProviderType;

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
                    this.idpType = IdpProviderType.IDPPROVIDERTYPE_ORG;
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.idpType = IdpProviderType.IDPPROVIDERTYPE_SYSTEM;
                    break;
            }
        }

        this.loadIdps();
    }

    public loadIdps(): void {
        this.idp = undefined;
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            const query: IdpSearchQuery = new IdpSearchQuery();
            query.setKey(IdpSearchKey.IDPSEARCHKEY_PROVIDER_TYPE);
            query.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
            query.setValue(this.idpType.toString());

            this.mgmtService.SearchIdps(undefined, undefined, [query]).then(idps => {
                this.availableIdps = idps.resultList;
            });
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
            this.adminService.listIDPs().then(idps => {
                this.availableIdps = idps.resultList;
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
