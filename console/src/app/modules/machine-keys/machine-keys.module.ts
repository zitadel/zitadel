import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
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
import { TableActionsModule } from '../table-actions/table-actions.module';
import { WarnDialogModule } from '../warn-dialog/warn-dialog.module';
import { MachineKeysComponent } from './machine-keys.component';

@NgModule({
  declarations: [MachineKeysComponent],
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
    TableActionsModule,
    WarnDialogModule,
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
  exports: [MachineKeysComponent],
})
export class MachineKeysModule {}
