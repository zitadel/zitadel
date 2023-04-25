import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatMomentDateModule } from '@angular/material-moment-adapter';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';

import { InfoSectionModule } from '../info-section/info-section.module';
import { AddTokenDialogComponent } from './add-token-dialog.component';

@NgModule({
  declarations: [AddTokenDialogComponent],
  imports: [
    CommonModule,
    TranslateModule,
    MatButtonModule,
    InfoSectionModule,
    InputModule,
    MatSelectModule,
    MatIconModule,
    FormsModule,
    MatDatepickerModule,
    MatMomentDateModule,
    ReactiveFormsModule,
    LocalizedDatePipeModule,
  ],
})
export class AddTokenDialogModule {}
