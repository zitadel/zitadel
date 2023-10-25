import { Pipe, PipeTransform } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { LocalizedDatePipe } from '../localized-date-pipe/localized-date.pipe';
import { TimestampToDatePipe } from '../timestamp-to-date-pipe/timestamp-to-date.pipe';
import { Milestone } from '../../proto/generated/zitadel/milestone/v1/milestone_pb';

@Pipe({
  name: 'milestone',
})
export class MilestonePipe implements PipeTransform {
  constructor(private translateService: TranslateService) {}

  public transform(milestone?: Milestone.AsObject): any {
    if (milestone && milestone.reachedDate) {
      const timestampToDate = new TimestampToDatePipe().transform(milestone.reachedDate);
      const datePipeOutput = new LocalizedDatePipe(this.translateService).transform(timestampToDate);
      return `done on ${datePipeOutput}`;
    } else {
      return '';
    }
  }
}
