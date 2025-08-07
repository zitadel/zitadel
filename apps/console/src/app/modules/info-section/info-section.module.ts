import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { InfoSectionComponent } from './info-section.component';

@NgModule({
  declarations: [InfoSectionComponent],
  imports: [CommonModule, TranslateModule, RouterModule],
  exports: [InfoSectionComponent],
})
export class InfoSectionModule {}
