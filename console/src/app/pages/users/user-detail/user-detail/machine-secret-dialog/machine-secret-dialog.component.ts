import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-machine-secret-dialog',
  templateUrl: './machine-secret-dialog.component.html',
  styleUrls: ['./machine-secret-dialog.component.scss'],
})
export class MachineSecretDialogComponent {
  public copied: string = '';
  constructor(
    public dialogRef: MatDialogRef<MachineSecretDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
