import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Idp, IdpProviderType, UserView } from 'src/app/proto/generated/management_pb';

@Component({
    selector: 'app-add-idp-dialog',
    templateUrl: './add-idp-dialog.component.html',
    styleUrls: ['./add-idp-dialog.component.scss'],
})
export class AddIdpDialogComponent {
    public preselectedIdps: Array<UserView.AsObject> = [];

    public idpType!: IdpProviderType;
    public idpTypes: IdpProviderType[] = [
        IdpProviderType.IDPPROVIDERTYPE_UNSPECIFIED,
        IdpProviderType.IDPPROVIDERTYPE_SYSTEM,
        IdpProviderType.IDPPROVIDERTYPE_ORG,
    ];

    public idps: Array<Idp.AsObject> | string[] = [];
    public IdpProviderType: any = IdpProviderType;

    constructor(
        public dialogRef: MatDialogRef<AddIdpDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
        if (data?.user) {
            this.preselectedIdps = [data.idps];
            this.idps = [data.idps];
        }
    }

    public loadIdps(): void {

    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({
            idps: this.idps,
        });
    }
}
