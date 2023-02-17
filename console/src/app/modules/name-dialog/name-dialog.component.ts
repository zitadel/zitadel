import { Component, Inject } from '@angular/core';
import {
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
  MatLegacyDialogRef as MatDialogRef,
} from '@angular/material/legacy-dialog';

@Component({
  selector: 'cnsl-name-dialog',
  templateUrl: './name-dialog.component.html',
  styleUrls: ['./name-dialog.component.scss'],
})
export class NameDialogComponent {
  public name: string = '';
  constructor(public dialogRef: MatDialogRef<NameDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.name = data.name ?? '';
  }

  closeDialog(name: string = ''): void {
    this.dialogRef.close(name);
  }
}
