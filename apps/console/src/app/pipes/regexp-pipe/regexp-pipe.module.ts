import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { RegexpPipe } from './regexp.pipe';

@NgModule({
  declarations: [RegexpPipe],
  imports: [CommonModule],
  exports: [RegexpPipe],
})
export class RegExpPipeModule {}
