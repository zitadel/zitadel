import { Component, Inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import {
  MAT_DIALOG_DATA,
  MatDialogModule,
  MatDialogRef,
} from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { InfoSectionType } from '../../info-section/info-section.component';

@Component({
  selector: 'cnsl-actions-two-signing-key-modal',
  templateUrl: './actions-two-signing-key-modal.component.html',
  styleUrls: ['./actions-two-signing-key-modal.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    TranslateModule,
    MatIconModule,
    MatTooltipModule,
    CopyToClipboardModule,
    InfoSectionModule,
  ],
})
export class ActionsTwoSigningKeyModalComponent {
  protected copied = '';
  protected readonly InfoSectionType = InfoSectionType;

  constructor(
    public dialogRef: MatDialogRef<ActionsTwoSigningKeyModalComponent>,
    @Inject(MAT_DIALOG_DATA) public readonly data: { signingKey: string },
  ) {}

  public close() {
    this.dialogRef.close();
  }
}
