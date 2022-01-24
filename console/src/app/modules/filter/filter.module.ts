import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { OutsideClickModule } from 'src/app/directives/outside-click/outside-click.module';

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
    MatCheckboxModule,
    InputModule,
    MatSelectModule,
    OutsideClickModule,
  ],
  exports: [FilterComponent],
})
export class FilterModule {}
