import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-resend-email-dialog',
    templateUrl: './resend-email-dialog.component.html',
    styleUrls: ['./resend-email-dialog.component.scss'],
})
export class ResendEmailDialogComponent {
    public email: string = '';
    constructor(public dialogRef: MatDialogRef<ResendEmailDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) { }

    closeDialog(email: string = ''): void {
        this.dialogRef.close(email);
    }

    closeDialogWithSend(email: string = ''): void {
        this.dialogRef.close({ send: true, email });
    }
}
