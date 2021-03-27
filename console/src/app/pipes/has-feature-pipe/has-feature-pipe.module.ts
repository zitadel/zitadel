import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { HasFeaturePipe } from './has-feature.pipe';


@NgModule({
    declarations: [
        HasFeaturePipe,
    ],
    imports: [
        CommonModule,
    ],
    exports: [
        HasFeaturePipe,
    ],
})
export class HasFeaturePipeModule { }
