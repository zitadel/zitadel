import { Component, inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

export type AppSecretDialogData = {
  clientId?: string;
  clientSecret?: string;
};

@Component({
  selector: 'cnsl-app-secret-dialog',
  templateUrl: './app-secret-dialog.component.html',
  styleUrls: ['./app-secret-dialog.component.scss'],
  standalone: false,
})
export class AppSecretDialogComponent {
  protected readonly dialogRef = inject<MatDialogRef<AppSecretDialogComponent>>(MatDialogRef);
  protected readonly data = inject<AppSecretDialogData>(MAT_DIALOG_DATA);

  public copied: string = '';
}
