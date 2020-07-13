import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
    selector: 'app-add-domain-dialog',
    templateUrl: './add-domain-dialog.component.html',
    styleUrls: ['./add-domain-dialog.component.scss'],
})
export class AddDomainDialogComponent implements OnInit {
    public newdomain: string = '';
    constructor(
        public dialogRef: MatDialogRef<AddDomainDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) { }

    ngOnInit(): void {
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close(this.newdomain);
    }
}
