import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { IdpView as AdminIdpView } from 'src/app/proto/generated/admin_pb';
import {
    Idp,
    IdpProviderType,
    IdpSearchKey,
    IdpSearchQuery,
    IdpView as MgmtIdpView,
    SearchMethod,
} from 'src/app/proto/generated/management_pb';
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
    public idpType!: IdpProviderType;
    public idpTypes: IdpProviderType[] = [
        IdpProviderType.IDPPROVIDERTYPE_SYSTEM,
        IdpProviderType.IDPPROVIDERTYPE_ORG,
    ];

    public idp: Idp.AsObject | undefined = undefined;
    public availableIdps: Array<AdminIdpView.AsObject | MgmtIdpView.AsObject> | string[] = [];
    public IdpProviderType: any = IdpProviderType;

    constructor(
        private mgmtService: ManagementService,
        private adminService: AdminService,
        public dialogRef: MatDialogRef<AddIdpDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
        if (data.idpType) {
            this.idpType = data.idpType;
            switch (this.idpType) {
                case IdpProviderType.IDPPROVIDERTYPE_ORG:
                    this.idpType = IdpProviderType.IDPPROVIDERTYPE_ORG;
                    break;
                case IdpProviderType.IDPPROVIDERTYPE_SYSTEM:
                    this.idpType = IdpProviderType.IDPPROVIDERTYPE_SYSTEM;
                    break;
            }
            this.loadIdps();
        }

    }

    public loadIdps(): void {
        console.log(this.idpType);
        this.idp = undefined;
        if (this.idpType === IdpProviderType.IDPPROVIDERTYPE_ORG) {
            const query: IdpSearchQuery = new IdpSearchQuery();
            query.setKey(IdpSearchKey.IDPSEARCHKEY_PROVIDER_TYPE);
            query.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
            query.setValue(this.idpType.toString());
            this.mgmtService.SearchIdps().then(idps => {
                console.log('mgmt', idps.toObject().resultList);
                this.availableIdps = idps.toObject().resultList;
            });
        } else if (this.idpType === IdpProviderType.IDPPROVIDERTYPE_SYSTEM) {
            this.adminService.SearchIdps().then(idps => {
                console.log('admin', idps.toObject().resultList);
                this.availableIdps = idps.toObject().resultList;
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
