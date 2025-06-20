import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { FilterModule } from '../filter/filter.module';
import { InputModule } from '../input/input.module';
import { FilterGroupComponent } from './filter-group.component';

@NgModule({
  declarations: [FilterGroupComponent],
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
  exports: [FilterGroupComponent],
})
export class FilterGroupModule {}
