import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { CreateActionRequest } from 'src/app/proto/generated/zitadel/management_pb';


@Component({
    selector: 'cnsl-add-action-dialog',
    templateUrl: './add-action-dialog.component.html',
    styleUrls: ['./add-action-dialog.component.scss'],
})
export class AddActionDialogComponent {
    public name: string = '';
    public script: string = '';
    public durationInSec: number = 10;
    public allowedToFail: boolean = false;
    
    constructor(
        public dialogRef: MatDialogRef<AddActionDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
       
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        const req = new CreateActionRequest();
        req.setName(this.name);
        req.setScript(this.script);

        const duration = new Duration();
        duration.setNanos(0);
        duration.setSeconds(this.durationInSec);

        req.setAllowedToFail(this.allowedToFail);

        req.setTimeout(duration)
        this.dialogRef.close(req);
    }
}
