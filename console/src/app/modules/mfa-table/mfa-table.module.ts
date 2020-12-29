import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { TruncatePipeModule } from 'src/app/pipes/truncate-pipe/truncate-pipe.module';

import { MfaTableComponent } from './mfa-table.component';
import { DialogAddTypeComponent } from './dialog-add-type/dialog-add-type.component';
import { InputModule } from '../input/input.module';
import { MatSelectModule } from '@angular/material/select';

@NgModule({
    declarations: [MfaTableComponent, DialogAddTypeComponent],
    imports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatIconModule,
        InputModule,
        MatSelectModule,
        MatTooltipModule,
        TranslateModule,
        TimestampToDatePipeModule,
        HasRoleModule,
        MatProgressSpinnerModule,
    ],
    exports: [
        MfaTableComponent,
    ],
})
export class MfaTableModule { }
