import {Component, EventEmitter, Inject, Input, Output} from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Observable } from 'rxjs';

@Component({
  templateUrl: './provider-next-dialog.component.html',
  styleUrls: ['./provider-next-dialog.component.scss'],
})
export class ProviderNextDialogComponent {
  @Output() activate = new EventEmitter<void>();
  constructor(
    public dialogRef: MatDialogRef<ProviderNextDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public next$: Observable<any>,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
