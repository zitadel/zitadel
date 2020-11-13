import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ErrorComponent } from './error.component';

@NgModule({
    declarations: [ErrorComponent],
    imports: [
        CommonModule,
    ],
    exports: [
        ErrorComponent,
    ],
})
export class ErrorModule { }

