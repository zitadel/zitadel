import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';

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
        MatRippleModule,
    ],
    exports: [InputDirective, FormFieldModule],
    providers: [ErrorStateMatcher],
})
export class InputModule { }
