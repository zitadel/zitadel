import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { TranslateModule } from '@ngx-translate/core';

import { InfoSectionModule } from '../info-section/info-section.module';
import { InputModule } from '../input/input.module';
import { WarnDialogComponent } from './warn-dialog.component';

@NgModule({
  declarations: [WarnDialogComponent],
  imports: [CommonModule, FormsModule, TranslateModule, InfoSectionModule, MatButtonModule, InputModule],
})
export class WarnDialogModule {}
