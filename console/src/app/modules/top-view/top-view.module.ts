import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { TopViewComponent } from './top-view.component';

@NgModule({
  declarations: [TopViewComponent],
  imports: [CommonModule, RouterModule, MatMenuModule, MatIconModule, MatButtonModule, MatTooltipModule, TranslateModule],
  exports: [TopViewComponent],
})
export class TopViewModule {}
