import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-dialog-u2f',
    templateUrl: './dialog-u2f.component.html',
    styleUrls: ['./dialog-u2f.component.scss'],
})
export class DialogU2FComponent {
    public name: string = '';
    constructor(public dialogRef: MatDialogRef<DialogU2FComponent>,
        @Inject(MAT_DIALOG_DATA) public data: string) { }

    public closeDialog(): void {
        this.dialogRef.close();
    }

    public closeDialogWithCode(): void {
        this.dialogRef.close(this.name);
    }
}
