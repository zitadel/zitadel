import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';

import { InputModule } from '../input/input.module';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatLegacyChipsModule } from '@angular/material/legacy-chips';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { AvatarModule } from '../avatar/avatar.module';
import { MatLegacyCheckboxModule } from '@angular/material/legacy-checkbox';
import { FilterEventsUserComponent } from './filter-events-user.component';

@NgModule({
  declarations: [FilterEventsUserComponent],
  imports: [
    CommonModule,
    MatAutocompleteModule,
    MatLegacyChipsModule,
    MatButtonModule,
    InputModule,
    MatIconModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    FormsModule,
    MatTooltipModule,
    TranslateModule,
    MatLegacyCheckboxModule,
    MatSelectModule,
    AvatarModule,
  ],
  exports: [FilterEventsUserComponent],
})
export class FilterEventsUserModule {}
