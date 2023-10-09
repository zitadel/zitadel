import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
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
