import { Component, Inject } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-add-token-dialog',
  templateUrl: './add-token-dialog.component.html',
  styleUrls: ['./add-token-dialog.component.scss'],
  standalone: false,
})
export class AddTokenDialogComponent {
  public startDate: Date = new Date();
  public dateControl: UntypedFormControl = new UntypedFormControl('', []);

  constructor(
    public dialogRef: MatDialogRef<AddTokenDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    const today = new Date();
    this.startDate.setDate(today.getDate() + 1);
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({ date: this.dateControl.value });
  }
}
