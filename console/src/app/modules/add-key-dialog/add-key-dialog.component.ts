import { Component, Inject } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { KeyType } from 'src/app/proto/generated/zitadel/auth_n_key_pb';

export enum AddKeyDialogType {
    MACHINE = 'MACHINE',
    AUTHNKEY = 'AUTHNKEY',
}

@Component({
    selector: 'app-add-key-dialog',
    templateUrl: './add-key-dialog.component.html',
    styleUrls: ['./add-key-dialog.component.scss'],
})
export class AddKeyDialogComponent {
    public startDate: Date = new Date();
    types: KeyType[] = [];
    public type!: KeyType;
    public dateControl: FormControl = new FormControl('', []);

    constructor(
        public dialogRef: MatDialogRef<AddKeyDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
        this.types = [KeyType.KEY_TYPE_JSON];
        this.type = KeyType.KEY_TYPE_JSON;
        const today = new Date();
        this.startDate.setDate(today.getDate() + 1);
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({ type: this.type, date: this.dateControl.value });
    }
}
