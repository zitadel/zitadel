import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';


@Component({
    selector: 'cnsl-add-action-dialog',
    templateUrl: './add-action-dialog.component.html',
    styleUrls: ['./add-action-dialog.component.scss'],
})
export class AddActionDialogComponent {
    public name: string = '';
    public script: string = '';
    
    constructor(
        public dialogRef: MatDialogRef<AddActionDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
       
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({  });
    }
}
