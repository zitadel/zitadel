import { Component, Inject } from '@angular/core';
import {
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
  MatLegacyDialogRef as MatDialogRef,
} from '@angular/material/legacy-dialog';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';

@Component({
  selector: 'cnsl-display-json-dialog',
  templateUrl: './display-json-dialog.component.html',
  styleUrls: ['./display-json-dialog.component.scss'],
})
export class DisplayJsonDialogComponent {
  public event?: Event.AsObject;
  public payload: any = '';

  constructor(public dialogRef: MatDialogRef<DisplayJsonDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.event = data.event;
    this.payload = data.event.payload.fieldsMap;
    console.log(this.event);
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
