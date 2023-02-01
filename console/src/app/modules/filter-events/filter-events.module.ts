import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';

import { InputModule } from '../input/input.module';
import { FilterEventsComponent } from './filter-events.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatLegacyProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyCheckboxModule } from '@angular/material/legacy-checkbox';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { OverlayModule } from '@angular/cdk/overlay';
import { MatDatepickerModule } from '@angular/material/datepicker';

@NgModule({
  declarations: [FilterEventsComponent],
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    InputModule,
    ReactiveFormsModule,
    OverlayModule,
    MatDatepickerModule,
    MatLegacyProgressSpinnerModule,
    TranslateModule,
    MatLegacyCheckboxModule,
    MatSelectModule,
    ActionKeysModule,
  ],
  exports: [FilterEventsComponent],
})
export class FilterEventsModule {}
