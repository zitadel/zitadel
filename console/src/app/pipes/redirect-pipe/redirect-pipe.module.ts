import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { RedirectPipe } from './redirect.pipe';


@NgModule({
    declarations: [
        RedirectPipe,
    ],
    imports: [
        CommonModule,
    ],
    exports: [
        RedirectPipe,
    ],
})
export class RedirectPipeModule { }
