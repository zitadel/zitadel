import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

import { MatDialogModule } from '@angular/material/dialog';
import { InfoSectionModule } from '../info-section/info-section.module';
import { InputModule } from '../input/input.module';
import { WarnDialogComponent } from './warn-dialog.component';

@NgModule({
  declarations: [WarnDialogComponent],
  imports: [CommonModule, FormsModule, MatDialogModule, TranslateModule, InfoSectionModule, MatButtonModule, InputModule],
})
export class WarnDialogModule {}
