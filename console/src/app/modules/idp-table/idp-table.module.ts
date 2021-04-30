import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { TruncatePipeModule } from 'src/app/pipes/truncate-pipe/truncate-pipe.module';

import { PaginatorModule } from '../paginator/paginator.module';
import { IdpTableComponent } from './idp-table.component';

@NgModule({
    declarations: [IdpTableComponent],
    imports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatCheckboxModule,
        MatIconModule,
        MatTooltipModule,
        TranslateModule,
        LocalizedDatePipeModule,
        TimestampToDatePipeModule,
        MatTableModule,
        PaginatorModule,
        RouterModule,
        RefreshTableModule,
        HasRoleModule,
        TruncatePipeModule,
    ],
    exports: [
        IdpTableComponent,
    ],
})
export class IdpTableModule { }
