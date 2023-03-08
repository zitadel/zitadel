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
import { FilterProjectComponent } from './filter-project.component';

@NgModule({
  declarations: [FilterProjectComponent],
  imports: [
    CommonModule,
    FilterModule,
    InputModule,
    MatSelectModule,
    MatCheckboxModule,
    MatButtonModule,
    MatIconModule,
    TranslateModule,
    RouterModule,
  ],
  exports: [FilterProjectComponent],
})
export class FilterProjectModule {}
