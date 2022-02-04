import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';

import { InfoOverlayComponent } from './info-overlay.component';

@NgModule({
  declarations: [InfoOverlayComponent],
  imports: [CommonModule, MatButtonModule, OverlayModule, MatIconModule, TranslateModule],
  exports: [InfoOverlayComponent],
})
export class InfoOverlayModule {}
