import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { FilterModule } from '../filter/filter.module';
import { InputModule } from '../input/input.module';
import { FilterUserComponent } from './filter-user.component';

@NgModule({
  declarations: [FilterUserComponent],
  imports: [
    CommonModule,
    FilterModule,
    InputModule,
    MatSelectModule,
    MatCheckboxModule,
    MatButtonModule,
    RouterModule,
    MatIconModule,
    TranslateModule,
  ],
  exports: [FilterUserComponent],
})
export class FilterUserModule {}
