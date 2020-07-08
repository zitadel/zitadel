import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-warn-dialog',
    templateUrl: './warn-dialog.component.html',
    styleUrls: ['./warn-dialog.component.scss'],
})
export class WarnDialogComponent implements OnInit {

    constructor(
        public dialogRef: MatDialogRef<WarnDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) { }

    ngOnInit(): void {
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close(true);
    }
}
