import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';

import { OverlayModule } from '@angular/cdk/overlay';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { InputModule } from '../input/input.module';
import { FilterEventsComponent } from './filter-events.component';
import {DateTimeLocalInputModule} from "../input-datetime-local/input-datetime-local.module";

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
    MatProgressSpinnerModule,
    TranslateModule,
    MatCheckboxModule,
    MatSelectModule,
    ActionKeysModule,
    DateTimeLocalInputModule,
  ],
  exports: [FilterEventsComponent],
})
export class FilterEventsModule {}
