import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-app-secret-dialog',
    templateUrl: './app-secret-dialog.component.html',
    styleUrls: ['./app-secret-dialog.component.scss'],
})
export class AppSecretDialogComponent {
    public copied: boolean = false;
    constructor(public dialogRef: MatDialogRef<AppSecretDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) { }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public copytoclipboard(value: string): void {
        const selBox = document.createElement('textarea');
        selBox.style.position = 'fixed';
        selBox.style.left = '0';
        selBox.style.top = '0';
        selBox.style.opacity = '0';
        selBox.value = value;
        document.body.appendChild(selBox);
        selBox.focus();
        selBox.select();
        document.execCommand('copy');
        document.body.removeChild(selBox);
        this.copied = true;
        setTimeout(() => {
            this.copied = false;
        }, 3000);
    }
}
