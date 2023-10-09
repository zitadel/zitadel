import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';

@Component({
  selector: 'cnsl-add-domain-dialog',
  templateUrl: './add-domain-dialog.component.html',
  styleUrls: ['./add-domain-dialog.component.scss'],
})
export class AddDomainDialogComponent {
  public newdomain: string = '';
  constructor(public dialogRef: MatDialogRef<AddDomainDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close(this.newdomain);
  }
}
