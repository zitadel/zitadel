import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { MultiFactor as AdminMultiFactor, MultiFactorType as AdminMultiFactorType } from 'src/app/proto/generated/admin_pb';
import { MultiFactor as MgmtMultiFactor, MultiFactorType as MgmtMultiFactorType } from 'src/app/proto/generated/management_pb';
import { LoginMethodComponentType } from '../mfa-table.component';

@Component({
    selector: 'app-dialog-add-type',
    templateUrl: './dialog-add-type.component.html',
    styleUrls: ['./dialog-add-type.component.scss'],
})
export class DialogAddTypeComponent {
    public LoginMethodComponentType: any = LoginMethodComponentType;
    public newMfaType!: AdminMultiFactorType | MgmtMultiFactorType;
    public availableMfaTypes: Array<AdminMultiFactorType | MgmtMultiFactorType> = [];
    constructor(public dialogRef: MatDialogRef<DialogAddTypeComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) {
        this.availableMfaTypes = data.types;
    }

    public closeDialog(): void {
        this.dialogRef.close();
    }

    public closeDialogWithCode(): void {
        this.dialogRef.close(this.newMfaType);
    }
}
