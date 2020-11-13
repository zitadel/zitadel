import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';
import { LabelModule } from 'src/app/modules/label/label.module';

import { ErrorComponent } from '../error/error.component';
import { ErrorModule } from '../error/error.module';
import { LabelComponent } from '../label/label.component';
import { FormFieldComponent } from './form-field.component';


@NgModule({
    declarations: [FormFieldComponent],
    imports: [
        CommonModule,
        MatRippleModule,
        LabelModule,
        ErrorModule,
    ],
    exports: [
        FormFieldComponent,
        LabelComponent,
        ErrorComponent,
    ],
})
export class FormFieldModule { }

