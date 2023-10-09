import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';

@Component({
  selector: 'cnsl-machine-secret-dialog',
  templateUrl: './machine-secret-dialog.component.html',
  styleUrls: ['./machine-secret-dialog.component.scss'],
})
export class MachineSecretDialogComponent {
  public copied: string = '';
  constructor(public dialogRef: MatDialogRef<MachineSecretDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
