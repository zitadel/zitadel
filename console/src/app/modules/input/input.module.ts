import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { FormFieldModule } from '../form-field/form-field.module';
import { LabelModule } from '../label/label.module';
import { ErrorStateMatcher } from './error-options';
import { InputDirective } from './input.directive';


@NgModule({
    declarations: [InputDirective],
    imports: [
        LabelModule,
        CommonModule,
        FormFieldModule,
    ],
    exports: [InputDirective, FormFieldModule],
    providers: [ErrorStateMatcher],
})
export class InputModule { }
