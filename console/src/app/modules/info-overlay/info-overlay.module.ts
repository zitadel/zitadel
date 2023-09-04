import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

import { InfoOverlayComponent } from './info-overlay.component';

@NgModule({
  declarations: [InfoOverlayComponent],
  imports: [CommonModule, MatButtonModule, MatIconModule, TranslateModule],
  exports: [InfoOverlayComponent],
})
export class InfoOverlayModule {}
