import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { AddKeyDialogModule } from 'src/app/modules/add-key-dialog/add-key-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { CardModule } from '../card/card.module';
import { InputModule } from '../input/input.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { ShowKeyDialogModule } from '../show-key-dialog/show-key-dialog.module';
import { ClientKeysComponent } from './client-keys.component';


@NgModule({
    declarations: [
        ClientKeysComponent,
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
        PaginatorModule,
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
        ClientKeysComponent,
    ],
})
export class ClientKeysModule { }
