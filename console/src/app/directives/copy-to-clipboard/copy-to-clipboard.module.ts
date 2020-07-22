import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { CopyToClipboardDirective } from './copy-to-clipboard.directive';



@NgModule({
    declarations: [
        CopyToClipboardDirective,
    ],
    imports: [
        CommonModule,
    ],
})
export class CopyToClipboardModule { }
