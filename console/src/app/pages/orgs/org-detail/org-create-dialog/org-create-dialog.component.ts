import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-org-create-dialog',
    templateUrl: './org-create-dialog.component.html',
    styleUrls: ['./org-create-dialog.component.scss'],
})
export class OrgCreateDialogComponent {
    constructor(
        private toast: ToastService,
        public dialogRef: MatDialogRef<OrgCreateDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        private orgService: OrgService,
    ) { }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }
}
