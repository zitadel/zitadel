import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-code-dialog',
  templateUrl: './code-dialog.component.html',
  styleUrls: ['./code-dialog.component.scss'],
})
export class CodeDialogComponent {
  public code: string = '';
  constructor(
    public dialogRef: MatDialogRef<CodeDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  closeDialog(code: string = ''): void {
    this.dialogRef.close(code);
  }
}
