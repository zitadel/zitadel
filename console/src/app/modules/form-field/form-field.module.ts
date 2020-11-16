import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { LabelModule } from 'src/app/modules/label/label.module';

import { LabelComponent } from '../label/label.component';
import { CnslErrorDirective } from './error.directive';
import { FormFieldComponent } from './form-field.component';


@NgModule({
    declarations: [FormFieldComponent, CnslErrorDirective],
    imports: [
        CommonModule,
        MatRippleModule,
        LabelModule,
    ],
    exports: [
        FormFieldComponent,
        LabelComponent,
        CnslErrorDirective,
    ],
})
export class FormFieldModule { }

