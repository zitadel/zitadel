import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-resend-email-dialog',
    templateUrl: './resend-email-dialog.component.html',
    styleUrls: ['./resend-email-dialog.component.scss'],
})
export class ResendEmailDialogComponent {
    public code: string = '';
    constructor(public dialogRef: MatDialogRef<ResendEmailDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) { }

    closeDialog(code: string = ''): void {
        this.dialogRef.close(code);
    }
}
