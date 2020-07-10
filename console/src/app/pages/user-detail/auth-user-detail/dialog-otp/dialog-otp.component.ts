import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-dialog-otp',
    templateUrl: './dialog-otp.component.html',
    styleUrls: ['./dialog-otp.component.scss'],
})
export class DialogOtpComponent {
    public code: string = '';
    constructor(public dialogRef: MatDialogRef<DialogOtpComponent>,
        @Inject(MAT_DIALOG_DATA) public data: string) { }

    public closeDialog(): void {
        this.dialogRef.close();
    }

    public closeDialogWithCode(): void {
        this.dialogRef.close(this.code);
    }
}
