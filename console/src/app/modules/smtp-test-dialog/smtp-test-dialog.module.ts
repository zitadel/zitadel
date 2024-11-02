import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

import { MatDialogModule } from '@angular/material/dialog';
import { InfoSectionModule } from '../info-section/info-section.module';
import { InputModule } from '../input/input.module';
import { SmtpTestDialogComponent } from './smtp-test-dialog.component';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@NgModule({
  declarations: [SmtpTestDialogComponent],
  imports: [
    CommonModule,
    FormsModule,
    MatDialogModule,
    MatProgressSpinnerModule,
    TranslateModule,
    InfoSectionModule,
    MatButtonModule,
    InputModule,
  ],
})
export class SmtpTestDialogModule {}
