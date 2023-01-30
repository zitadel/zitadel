import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';

import { InputModule } from '../input/input.module';
import { FilterEventsComponent } from './filter-events.component';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatLegacyChipsModule } from '@angular/material/legacy-chips';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatLegacyProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule } from '@angular/material/legacy-tooltip';
import { AvatarModule } from '../avatar/avatar.module';
import { MatLegacyCheckboxModule } from '@angular/material/legacy-checkbox';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { OverlayModule } from '@angular/cdk/overlay';

@NgModule({
  declarations: [FilterEventsComponent],
  imports: [
    CommonModule,
    FormsModule,
    MatAutocompleteModule,
    MatLegacyChipsModule,
    MatButtonModule,
    InputModule,
    MatIconModule,
    ReactiveFormsModule,
    OverlayModule,
    MatLegacyProgressSpinnerModule,
    MatLegacyTooltipModule,
    TranslateModule,
    MatLegacyCheckboxModule,
    MatSelectModule,
    AvatarModule,
    ActionKeysModule,
  ],
  exports: [FilterEventsComponent],
})
export class FilterEventsModule {}
