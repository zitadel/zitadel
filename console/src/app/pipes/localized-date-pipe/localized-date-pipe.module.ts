import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { LocalizedDatePipe } from './localized-date.pipe';


@NgModule({
    declarations: [
        LocalizedDatePipe,
    ],
    imports: [
        CommonModule,
    ],
    exports: [
        LocalizedDatePipe,
    ],
})
export class LocalizedDatePipeModule { }
