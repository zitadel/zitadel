import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';

import { InfoSectionType } from '../info-section/info-section.component';

@Component({
  selector: 'cnsl-warn-dialog',
  templateUrl: './warn-dialog.component.html',
  styleUrls: ['./warn-dialog.component.scss'],
})
export class WarnDialogComponent {
  public confirm: string = '';
  InfoSectionType: any = InfoSectionType;
  constructor(public dialogRef: MatDialogRef<WarnDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close(true);
  }
}
