import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MomentModule } from 'ngx-moment';

import { LocalizedDatePipe } from './localized-date.pipe';


@NgModule({
    declarations: [
        LocalizedDatePipe,
    ],
    imports: [
        CommonModule,
        MomentModule,
    ],
    exports: [
        LocalizedDatePipe,
    ],
})
export class PipesModule { }
