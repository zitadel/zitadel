import { Component, inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';

import { InfoSectionType } from '../info-section/info-section.component';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';

export type InfoDialogData = {
  confirmKey: string;
  cancelKey: string;
  titleKey: string;
  descriptionKey: string;
};

export type InfoDialogResult = boolean;

@Component({
  selector: 'cnsl-info-dialog',
  templateUrl: './info-dialog.component.html',
  styleUrls: ['./info-dialog.component.scss'],
  imports: [TranslateModule, MatDialogModule, MatButtonModule],
})
export class InfoDialogComponent {
  protected readonly InfoSectionType = InfoSectionType;

  protected readonly dialogRef = inject<MatDialogRef<InfoDialogComponent>>(MatDialogRef);
  protected readonly data = inject<InfoDialogData>(MAT_DIALOG_DATA);
}
