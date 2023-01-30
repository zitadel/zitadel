import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ToPayloadPipe } from './to-payload.pipe';

@NgModule({
  declarations: [ToPayloadPipe],
  imports: [CommonModule],
  exports: [ToPayloadPipe],
})
export class ToPayloadPipeModule {}
