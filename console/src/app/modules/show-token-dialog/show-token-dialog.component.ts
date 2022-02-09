import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { AddPersonalAccessTokenResponse } from 'src/app/proto/generated/zitadel/management_pb';

import { InfoSectionType } from '../info-section/info-section.component';

@Component({
  selector: 'cnsl-show-token-dialog',
  templateUrl: './show-token-dialog.component.html',
  styleUrls: ['./show-token-dialog.component.scss'],
})
export class ShowTokenDialogComponent {
  public tokenResponse!: AddPersonalAccessTokenResponse.AsObject;
  InfoSectionType: any = InfoSectionType;

  constructor(public dialogRef: MatDialogRef<ShowTokenDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.tokenResponse = data.key;
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
