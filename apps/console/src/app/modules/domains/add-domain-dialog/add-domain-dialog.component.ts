import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-add-domain-dialog',
  templateUrl: './add-domain-dialog.component.html',
  styleUrls: ['./add-domain-dialog.component.scss'],
})
export class AddDomainDialogComponent {
  public newdomain: string = '';
  constructor(
    public dialogRef: MatDialogRef<AddDomainDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close(this.newdomain);
  }
}
