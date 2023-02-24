import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { CreateLayoutComponent } from './create-layout.component';

@NgModule({
  declarations: [CreateLayoutComponent],
  imports: [CommonModule, MatIconModule, MatButtonModule, TranslateModule, MatTooltipModule],
  exports: [CreateLayoutComponent],
})
export class CreateLayoutModule {}
