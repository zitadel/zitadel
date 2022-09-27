import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { TimestampToDatePipe } from './timestamp-to-date.pipe';

@NgModule({
  declarations: [TimestampToDatePipe],
  imports: [CommonModule],
  exports: [TimestampToDatePipe],
})
export class TimestampToDatePipeModule {}
