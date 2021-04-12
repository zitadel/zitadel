import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { RefreshTableComponent } from './refresh-table.component';

@NgModule({
    declarations: [RefreshTableComponent],
    imports: [
        CommonModule,
        MatButtonModule,
        MatIconModule,
        TranslateModule,
        FormsModule,
        MatTooltipModule,
        MatProgressSpinnerModule,
        TimestampToDatePipeModule,
        LocalizedDatePipeModule,
        PaginatorModule,
    ],
    exports: [
        RefreshTableComponent,
    ],
})
export class RefreshTableModule { }
