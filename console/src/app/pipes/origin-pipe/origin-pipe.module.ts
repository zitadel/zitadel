import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { OriginPipe } from './origin.pipe';


@NgModule({
    declarations: [
        OriginPipe,
    ],
    imports: [
        CommonModule,
    ],
    exports: [
        OriginPipe,
    ],
})
export class OriginPipeModule { }
