import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ScrollableDirective } from './scrollable.directive';

@NgModule({
  declarations: [ScrollableDirective],
  imports: [CommonModule],
  exports: [ScrollableDirective],
})
export class ScrollableModule {}
