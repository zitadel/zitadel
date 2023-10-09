import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { mapTo } from 'rxjs';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';

@Component({
  selector: 'cnsl-display-json-dialog',
  templateUrl: './display-json-dialog.component.html',
  styleUrls: ['./display-json-dialog.component.scss'],
})
export class DisplayJsonDialogComponent {
  public event?: Event;
  public payload: any = '';
  public opened$ = this.dialogRef.afterOpened().pipe(mapTo(true));

  constructor(public dialogRef: MatDialogRef<DisplayJsonDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.event = data.event;
    if ((data.event as Event) && data.event.payload) {
    }
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
