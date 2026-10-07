import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

import { TableSearchComponent } from './table-search.component';

@NgModule({
  declarations: [TableSearchComponent],
  imports: [CommonModule, TranslateModule],
  exports: [TableSearchComponent],
})
export class TableSearchModule {}
