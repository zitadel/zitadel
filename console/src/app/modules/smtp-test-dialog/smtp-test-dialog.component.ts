import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { InfoSectionType } from '../info-section/info-section.component';

@Component({
  selector: 'cnsl-smtp-test-dialog',
  templateUrl: './smtp-test-dialog.component.html',
  styleUrls: ['./smtp-test-dialog.component.scss'],
})
export class SmtpTestDialogComponent {
  public confirm: string = '';
  InfoSectionType: any = InfoSectionType;
  constructor(
    public dialogRef: MatDialogRef<SmtpTestDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close(true);
  }
}
