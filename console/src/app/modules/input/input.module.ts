import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { InputDirective } from './input.directive';



@NgModule({
    declarations: [InputDirective],
    imports: [
        CommonModule,
    ],
    exports: [InputDirective],
})
export class InputModule { }
