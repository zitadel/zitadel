import { DatePipe } from '@angular/common';
import { Pipe, PipeTransform } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import * as moment from 'moment';

@Pipe({
  name: 'localizedDate',
})
export class LocalizedDatePipe implements PipeTransform {
  constructor(private translateService: TranslateService) {}

  public transform(value: any, pattern?: string): any {
    if (pattern && pattern === 'fromNow') {
      moment.locale(this.translateService.currentLang ?? 'en');

      let date = moment(value);
      if (moment().diff(date, 'days') <= 2) {
        return date.fromNow(); // '2 days ago' etc.
      } else {
        const localeData = moment(value).localeData();
        const format = localeData.longDateFormat('L');
        return moment(value).format(`${format}, HH:mm`);
      }
    } else {
      const lang = ['de', 'en', 'fr', 'it', 'pl', 'zh'].includes(this.translateService.currentLang)
        ? this.translateService.currentLang
        : 'en';
      const datePipe: DatePipe = new DatePipe(lang);
      return datePipe.transform(value, pattern ?? 'mediumDate');
    }
  }
}
