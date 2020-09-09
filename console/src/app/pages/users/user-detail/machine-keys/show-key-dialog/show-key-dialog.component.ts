import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { AddMachineKeyResponse } from 'src/app/proto/generated/management_pb';

@Component({
    selector: 'app-show-key-dialog',
    templateUrl: './show-key-dialog.component.html',
    styleUrls: ['./show-key-dialog.component.scss'],
})
export class ShowKeyDialogComponent {
    public addedKey!: AddMachineKeyResponse.AsObject;

    constructor(
        public dialogRef: MatDialogRef<ShowKeyDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
        this.addedKey = data.key;
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }
}
