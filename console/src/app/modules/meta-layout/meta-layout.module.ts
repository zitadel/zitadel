import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { MetaLayoutComponent } from './meta-layout.component';



@NgModule({
    declarations: [MetaLayoutComponent],
    imports: [
        CommonModule,
    ],
    exports: [MetaLayoutComponent],
})
export class MetaLayoutModule { }
