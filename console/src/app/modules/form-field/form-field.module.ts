import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { LabelModule } from 'src/app/modules/label/label.module';

import { LabelComponent } from '../label/label.component';
import { CnslErrorDirective } from './error.directive';
import { CnslFormFieldComponent } from './form-field.component';
import { CnslHintDirective } from './hint.directive';

@NgModule({
    declarations: [
        CnslFormFieldComponent,
        CnslErrorDirective,
        CnslHintDirective,
    ],
    imports: [
        CommonModule,
        MatRippleModule,
        LabelModule,
    ],
    exports: [
        CnslFormFieldComponent,
        LabelComponent,
        CnslErrorDirective,
        CnslHintDirective,
    ],
})
export class FormFieldModule { }

