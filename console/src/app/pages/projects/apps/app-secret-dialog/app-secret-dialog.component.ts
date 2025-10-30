import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';

@Component({
  selector: 'cnsl-app-secret-dialog',
  templateUrl: './app-secret-dialog.component.html',
  styleUrls: ['./app-secret-dialog.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatIconModule,
    MatTooltipModule,
    TranslateModule,
    CopyToClipboardModule,
    InfoSectionModule,
  ],
})
export class AppSecretDialogComponent {
  public copied: string = '';
  constructor(
    public dialogRef: MatDialogRef<AppSecretDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
