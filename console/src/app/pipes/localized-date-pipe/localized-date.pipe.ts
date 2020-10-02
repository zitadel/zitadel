import { DatePipe } from '@angular/common';
import { Pipe, PipeTransform } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

@Pipe({
    name: 'localizedDate',
})
export class LocalizedDatePipe implements PipeTransform {

    constructor(private translateService: TranslateService) { }

    public transform(value: any, pattern: string = 'mediumDate'): any {
        if (this.translateService.currentLang && this.translateService.currentLang === ('de' || 'it' || 'fr' || 'eng')) {
            const datePipe: DatePipe = new DatePipe(this.translateService.currentLang);
            return datePipe.transform(value, pattern);
        } else {
            const datePipe: DatePipe = new DatePipe('de');
            return datePipe.transform(value, pattern);
        }
    }
}
