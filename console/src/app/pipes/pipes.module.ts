import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MomentModule } from 'ngx-moment';

import { LocalizedDatePipe } from './localized-date.pipe';
import { TimestampToDatePipe } from './timestamp-to-date.pipe';


@NgModule({
    declarations: [
        LocalizedDatePipe,
        TimestampToDatePipe,
    ],
    imports: [
        CommonModule,
        MomentModule,
    ],
    exports: [
        LocalizedDatePipe,
        TimestampToDatePipe,
    ],
})
export class PipesModule { }
