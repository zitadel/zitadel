import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatRippleModule } from '@angular/material/core';

import { LabelComponent } from './label.component';



@NgModule({
    declarations: [LabelComponent],
    imports: [
        CommonModule,
        MatRippleModule,
    ],
    exports: [
        LabelComponent,
    ],
})
export class LabelModule { }

