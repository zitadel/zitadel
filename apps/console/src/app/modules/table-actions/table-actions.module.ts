import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { TableActionsComponent } from './table-actions.component';

@NgModule({
  declarations: [TableActionsComponent],
  imports: [CommonModule, MatIconModule, MatButtonModule, MatMenuModule, MatTooltipModule, TranslateModule],
  exports: [TableActionsComponent, MatMenuModule],
})
export class TableActionsModule {}
