import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { saveAs } from 'file-saver';
import { AddMachineTokenResponse } from 'src/app/proto/generated/zitadel/management_pb';

@Component({
  selector: 'cnsl-show-token-dialog',
  templateUrl: './show-token-dialog.component.html',
  styleUrls: ['./show-token-dialog.component.scss'],
})
export class ShowTokenDialogComponent {
  public tokenResponse!: AddMachineTokenResponse.AsObject;

  constructor(public dialogRef: MatDialogRef<ShowTokenDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.tokenResponse = data.key;
  }

  public saveFile(): void {
    const json = atob(this.tokenResponse.token.toString());
    const blob = new Blob([json], { type: 'text/plain;charset=utf-8' });
    const name = this.tokenResponse.tokenId;
    saveAs(blob, `${name}.json`);
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
