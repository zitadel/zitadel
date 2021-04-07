import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { TimestampToRetentionPipe } from './timestamp-to-retention.pipe';


@NgModule({
    declarations: [
        TimestampToRetentionPipe,
    ],
    imports: [
        CommonModule,
    ],
    exports: [
        TimestampToRetentionPipe,
    ],
})
export class TimestampToRetentionPipeModule { }
