import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MomentModule } from 'ngx-moment';

import { LocalizedDatePipe } from './localized-date.pipe';
import { PasswordPatternPipe } from './password-pattern.pipe';


@NgModule({
    declarations: [
        LocalizedDatePipe,
        PasswordPatternPipe,
    ],
    imports: [
        CommonModule,
        MomentModule,
    ],
    exports: [
        LocalizedDatePipe,
        PasswordPatternPipe,
    ],
})
export class PipesModule { }
