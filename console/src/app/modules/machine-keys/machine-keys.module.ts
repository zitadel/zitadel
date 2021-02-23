import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { CardModule } from '../card/card.module';
import { InputModule } from '../input/input.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';

import { MachineKeysComponent } from './machine-keys.component';
import { ShowKeyDialogModule } from './show-key-dialog/show-key-dialog.module';
import { AddKeyDialogModule } from './add-key-dialog/add-key-dialog.module';
import { RouterModule } from '@angular/router';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';


@NgModule({
    declarations: [
        MachineKeysComponent,
    ],
    imports: [
        CommonModule,
        RouterModule,
        FormsModule,
        MatButtonModule,
        MatDialogModule,
        HasRoleModule,
        CardModule,
        MatTableModule,
        MatPaginatorModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatCheckboxModule,
        MatTooltipModule,
        HasRolePipeModule,
        TimestampToDatePipeModule,
        LocalizedDatePipeModule,
        TranslateModule,
        RefreshTableModule,
        InputModule,
        ShowKeyDialogModule,
        AddKeyDialogModule,
    ],
    exports: [
        MachineKeysComponent,
    ],
})
export class MachineKeysModule { }
