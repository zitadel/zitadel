import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';

import { OverlayModule } from '@angular/cdk/overlay';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatLegacyCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { InputModule } from '../input/input.module';
import { FilterEventsComponent } from './filter-events.component';

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
