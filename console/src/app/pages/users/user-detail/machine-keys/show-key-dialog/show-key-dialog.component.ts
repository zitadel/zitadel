import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { saveAs } from 'file-saver';
import { AddMachineKeyResponse } from 'src/app/proto/generated/zitadel/management_pb';

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

    public saveFile(): void {
        const json = atob(this.addedKey.keyDetails.toString());
        const blob = new Blob([json], { type: 'text/plain;charset=utf-8' });
        saveAs(blob, `${this.addedKey.id}.json`);
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }
}
