import { A11yModule } from '@angular/cdk/a11y';
import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';

import { ActionKeysModule } from '../action-keys/action-keys.module';
import { InputModule } from '../input/input.module';
import { FilterComponent } from './filter.component';

@NgModule({
  declarations: [FilterComponent],
  imports: [
    CommonModule,
    TranslateModule,
    ActionKeysModule,
    MatButtonModule,
    MatIconModule,
    OverlayModule,
    A11yModule,
    MatCheckboxModule,
    InputModule,
    MatSelectModule,
  ],
  exports: [FilterComponent],
})
export class FilterModule {}
