import { Component, Inject } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MachineKeyType } from 'src/app/proto/generated/management_pb';

@Component({
    selector: 'app-add-key-dialog',
    templateUrl: './add-key-dialog.component.html',
    styleUrls: ['./add-key-dialog.component.scss'],
})
export class AddKeyDialogComponent {
    public startDate: Date = new Date();
    types: MachineKeyType[] = [
        MachineKeyType.MACHINEKEY_JSON,
    ];
    public type: MachineKeyType = MachineKeyType.MACHINEKEY_JSON;
    public dateControl: FormControl = new FormControl('', []);

    constructor(
        public dialogRef: MatDialogRef<AddKeyDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
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
