import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ToObjectPipe } from './to-object.pipe';

@NgModule({
  declarations: [ToObjectPipe],
  imports: [CommonModule],
  exports: [ToObjectPipe],
})
export class ToObjectPipeModule {}
