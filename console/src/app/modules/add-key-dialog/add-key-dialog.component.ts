import { Component, Inject } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { KeyType } from 'src/app/proto/generated/zitadel/auth_n_key_pb';

export enum AddKeyDialogType {
  MACHINE = 'MACHINE',
  AUTHNKEY = 'AUTHNKEY',
}

@Component({
  selector: 'cnsl-add-key-dialog',
  templateUrl: './add-key-dialog.component.html',
  styleUrls: ['./add-key-dialog.component.scss'],
})
export class AddKeyDialogComponent {
  public startDate: Date = new Date();
  types: KeyType[] = [];
  public type!: KeyType;
  public dateControl: UntypedFormControl = new UntypedFormControl('', []);

  constructor(public dialogRef: MatDialogRef<AddKeyDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.types = [KeyType.KEY_TYPE_JSON];
    this.type = KeyType.KEY_TYPE_JSON;
    const today = new Date();
    this.startDate.setDate(today.getDate() + 1);
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({ type: this.type, date: this.dateControl.value });
  }
}
