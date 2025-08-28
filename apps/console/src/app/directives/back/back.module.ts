import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { BackDirective } from './back.directive';

@NgModule({
  declarations: [BackDirective],
  imports: [CommonModule],
  exports: [BackDirective],
})
export class BackModule {}
