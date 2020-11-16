import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ErrorDirective } from './error.directive';


@NgModule({
    declarations: [ErrorDirective],
    imports: [
        CommonModule,
    ],
    exports: [
        ErrorDirective,
    ],
})
export class ErrorModule { }

