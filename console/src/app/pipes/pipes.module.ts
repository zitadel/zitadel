import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MomentModule } from 'ngx-moment';

import { HasRolePipe } from './has-role.pipe';
import { LocalizedDatePipe } from './localized-date.pipe';
import { TimestampToDatePipe } from './timestamp-to-date.pipe';


@NgModule({
    declarations: [
        LocalizedDatePipe,
        TimestampToDatePipe,
        HasRolePipe,
    ],
    imports: [
        CommonModule,
        MomentModule,
    ],
    exports: [
        LocalizedDatePipe,
        TimestampToDatePipe,
        HasRolePipe,
    ],
})
export class PipesModule { }
