import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { DurationToSecondsPipe } from './duration-to-seconds.pipe';

@NgModule({
  declarations: [DurationToSecondsPipe],
  imports: [CommonModule],
  exports: [DurationToSecondsPipe],
})
export class DurationToSecondsPipeModule {}
