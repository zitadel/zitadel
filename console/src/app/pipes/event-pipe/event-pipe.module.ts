import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { LocalizedDatePipeModule } from '../localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from '../timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { EventPipe } from './event.pipe';

@NgModule({
  declarations: [EventPipe],
  imports: [CommonModule, TimestampToDatePipeModule, LocalizedDatePipeModule],
  exports: [EventPipe],
})
export class EventPipeModule {}
