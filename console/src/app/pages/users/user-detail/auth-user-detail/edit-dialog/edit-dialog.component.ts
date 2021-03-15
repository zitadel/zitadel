import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-edit-email-dialog',
    templateUrl: './edit-dialog.component.html',
    styleUrls: ['./edit-dialog.component.scss'],
})
export class EditDialogComponent {
    public value: string = '';
    constructor(public dialogRef: MatDialogRef<EditDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) {
        this.value = data.value;
        console.log(this.value);
    }

    closeDialog(email: string = ''): void {
        this.dialogRef.close(email);
    }

    closeDialogWithValue(value: string = ''): void {
        this.dialogRef.close(value);
    }
}
