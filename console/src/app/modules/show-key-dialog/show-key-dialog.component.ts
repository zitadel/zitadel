import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { saveAs } from 'file-saver';
import { AddAppKeyResponse, AddMachineKeyResponse } from 'src/app/proto/generated/zitadel/management_pb';

@Component({
  selector: 'app-show-key-dialog',
  templateUrl: './show-key-dialog.component.html',
  styleUrls: ['./show-key-dialog.component.scss'],
})
export class ShowKeyDialogComponent {
  public keyResponse!: AddMachineKeyResponse.AsObject | AddAppKeyResponse.AsObject;

  constructor(
    public dialogRef: MatDialogRef<ShowKeyDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.keyResponse = data.key;
  }

  public saveFile(): void {
    const json = atob(this.keyResponse.keyDetails.toString());
    const blob = new Blob([json], { type: 'text/plain;charset=utf-8' });
    const name = (this.keyResponse as AddMachineKeyResponse.AsObject).keyId ? (this.keyResponse as AddMachineKeyResponse.AsObject).keyId : (this.keyResponse as AddAppKeyResponse.AsObject).id;
    saveAs(blob, `${name}.json`);
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
