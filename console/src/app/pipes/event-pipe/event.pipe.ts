import { Pipe, PipeTransform } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { LocalizedDatePipe } from '../localized-date-pipe/localized-date.pipe';
import { TimestampToDatePipe } from '../timestamp-to-date-pipe/timestamp-to-date.pipe';

@Pipe({
  name: 'event',
})
export class EventPipe implements PipeTransform {
  constructor(private translateService: TranslateService) {}

  public transform(event?: Event.AsObject): any {
    if (event && event.editor?.displayName && event.creationDate) {
      const timestampToDate = new TimestampToDatePipe().transform(event.creationDate);
      const datePipeOutput = new LocalizedDatePipe(this.translateService).transform(timestampToDate);
      return `${event.editor?.displayName} last changed it on ${datePipeOutput}`;
    } else if (event && event.creationDate) {
      const timestampToDate = new TimestampToDatePipe().transform(event.creationDate);
      const datePipeOutput = new LocalizedDatePipe(this.translateService).transform(timestampToDate);
      return `done on ${datePipeOutput}`;
    } else {
      return '';
    }
  }
}
