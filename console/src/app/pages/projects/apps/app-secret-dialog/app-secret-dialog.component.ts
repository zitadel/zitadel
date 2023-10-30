import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-app-secret-dialog',
  templateUrl: './app-secret-dialog.component.html',
  styleUrls: ['./app-secret-dialog.component.scss'],
})
export class AppSecretDialogComponent {
  public copied: string = '';
  constructor(
    public dialogRef: MatDialogRef<AppSecretDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
