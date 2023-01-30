import { Component, Inject } from '@angular/core';
import {
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
  MatLegacyDialogRef as MatDialogRef,
} from '@angular/material/legacy-dialog';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { Struct } from 'google-protobuf/google/protobuf/struct_pb';

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
    if ((data.event as Event.AsObject) && data.event.payload) {
      console.log((data.event as Event.AsObject).payload?.fieldsMap);
      this.payload = (data.event as Event.AsObject).payload?.fieldsMap;
    }
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
