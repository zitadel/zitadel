import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialog as MatDialog,
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { mapTo } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Action } from 'src/app/proto/generated/zitadel/action_pb';
import { CreateActionRequest, UpdateActionRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

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

  public id: string = '';

  public opened$ = this.dialogRef.afterOpened().pipe(mapTo(true));

  constructor(
    private toast: ToastService,
    private mgmtService: ManagementService,
    private dialog: MatDialog,
    public dialogRef: MatDialogRef<AddActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    if (data && data.action) {
      const action: Action.AsObject = data.action;
      this.name = action.name;
      this.script = action.script;
      if (action.timeout?.seconds) {
        this.durationInSec = action.timeout?.seconds;
      }
      this.allowedToFail = action.allowedToFail;
      this.id = action.id;
    }
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    if (this.id) {
      const req = new UpdateActionRequest();
      req.setId(this.id);
      req.setName(this.name);
      req.setScript(this.script);

      const duration = new Duration();
      duration.setNanos(0);
      duration.setSeconds(this.durationInSec);

      req.setAllowedToFail(this.allowedToFail);

      req.setTimeout(duration);
      this.dialogRef.close(req);
    } else {
      const req = new CreateActionRequest();
      req.setName(this.name);
      req.setScript(this.script);

      const duration = new Duration();
      duration.setNanos(0);
      duration.setSeconds(this.durationInSec);

      req.setAllowedToFail(this.allowedToFail);

      req.setTimeout(duration);
      this.dialogRef.close(req);
    }
  }

  public deleteAndCloseDialog(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'FLOWS.DIALOG.DELETEACTION.TITLE',
        descriptionKey: 'FLOWS.DIALOG.DELETEACTION.DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .deleteAction(this.id)
          .then((resp) => {
            this.dialogRef.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });
  }
}
