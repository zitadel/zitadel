import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { TruncatePipePipe } from './truncate-pipe.pipe';

@NgModule({
  declarations: [TruncatePipePipe],
  imports: [CommonModule],
  exports: [TruncatePipePipe],
})
export class TruncatePipeModule {}
