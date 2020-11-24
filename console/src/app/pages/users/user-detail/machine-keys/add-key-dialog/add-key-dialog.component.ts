import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MachineKeyType } from 'src/app/proto/generated/management_pb';

@Component({
    selector: 'app-add-key-dialog',
    templateUrl: './add-key-dialog.component.html',
    styleUrls: ['./add-key-dialog.component.scss'],
})
export class AddKeyDialogComponent {
    types: MachineKeyType[] = [
        MachineKeyType.MACHINEKEY_JSON,
    ];
    date!: Date;
    public type: MachineKeyType = MachineKeyType.MACHINEKEY_JSON;

    constructor(
        public dialogRef: MatDialogRef<AddKeyDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) { }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({ type: this.type, date: this.date });
    }
}
