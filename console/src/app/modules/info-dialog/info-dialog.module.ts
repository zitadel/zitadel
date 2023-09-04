import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

import { InfoSectionModule } from '../info-section/info-section.module';
import { InputModule } from '../input/input.module';
import { InfoDialogComponent } from './info-dialog.component';

@NgModule({
  declarations: [InfoDialogComponent],
  imports: [CommonModule, FormsModule, TranslateModule, InfoSectionModule, MatButtonModule, InputModule],
})
export class InfoDialogModule {}
