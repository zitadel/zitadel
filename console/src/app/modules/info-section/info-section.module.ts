import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { InfoSectionComponent } from './info-section.component';



@NgModule({
    declarations: [InfoSectionComponent],
    imports: [
        CommonModule,
    ],
    exports: [
        InfoSectionComponent,
    ],
})
export class InfoSectionModule { }
