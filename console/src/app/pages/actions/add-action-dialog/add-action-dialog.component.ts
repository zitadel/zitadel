import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { mapTo, Subject, takeUntil } from 'rxjs';
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
export class AddActionDialogComponent implements OnInit, OnDestroy {
  public id: string = '';

  public opened$ = this.dialogRef.afterOpened().pipe(mapTo(true));
  private destroy$: Subject<void> = new Subject();
  private showAllowedToFailWarning: boolean = true;

  public form: FormGroup = new FormGroup({
    name: new FormControl<string>('', []),
    script: new FormControl<string>('', []),
    durationInSec: new FormControl<number>(10, []),
    allowedToFail: new FormControl<boolean>(true, []),
  });
  constructor(
    private toast: ToastService,
    private mgmtService: ManagementService,
    private dialog: MatDialog,
    private unsavedChangesDialog: MatDialog,
    public dialogRef: MatDialogRef<AddActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    if (data && data.action) {
      const action: Action.AsObject = data.action;
      this.form.setValue({
        name: action.name,
        script: action.script,
        durationInSec: action.timeout?.seconds ?? 10,
        allowedToFail: action.allowedToFail,
      });
      this.id = action.id;
    }

    this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe(({ allowedToFail }) => {
      if (!allowedToFail && this.showAllowedToFailWarning) {
        this.dialog.open(WarnDialogComponent, {
          data: {
            confirmKey: 'ACTIONS.OK',
            titleKey: 'FLOWS.ALLOWEDTOFAILWARN.TITLE',
            descriptionKey: 'FLOWS.ALLOWEDTOFAILWARN.DESCRIPTION',
          },
          width: '400px',
        });

        this.showAllowedToFailWarning = false;
      }
    });
  }

  ngOnInit(): void {
    // prevent unsaved changes get lost if backdrop is clicked
    this.dialogRef.backdropClick().subscribe(() => {
      if (this.form.dirty) {
        this.showUnsavedDialog();
      } else {
        this.dialogRef.close(false);
      }
    });

    // prevent unsaved changes get lost if escape key is pressed
    this.dialogRef.keydownEvents().subscribe((event) => {
      if (event.key === 'Escape') {
        if (this.form.dirty) {
          this.showUnsavedDialog();
        } else {
          this.dialogRef.close(false);
        }
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private showUnsavedDialog(): void {
    const unsavedChangesDialogRef = this.unsavedChangesDialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.UNSAVED.DIALOG.DISCARD',
        cancelKey: 'ACTIONS.UNSAVED.DIALOG.CANCEL',
        titleKey: 'ACTIONS.UNSAVEDCHANGES',
        descriptionKey: 'ACTIONS.UNSAVED.DIALOG.DESCRIPTION',
      },
      width: '400px',
    });

    unsavedChangesDialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.dialogRef.close(false);
      }
    });
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    if (this.id) {
      const req = new UpdateActionRequest();
      req.setId(this.id);
      req.setName(this.form.value.name);
      req.setScript(this.form.value.script);

      const duration = new Duration();
      duration.setNanos(0);
      duration.setSeconds(this.form.value.durationInSec);

      req.setAllowedToFail(this.form.value.allowedToFail);

      req.setTimeout(duration);
      this.dialogRef.close(req);
    } else {
      const req = new CreateActionRequest();
      req.setName(this.form.value.name);
      req.setScript(this.form.value.script);

      const duration = new Duration();
      duration.setNanos(0);
      duration.setSeconds(this.form.value.durationInSec);

      req.setAllowedToFail(this.form.value.allowedToFail);

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
