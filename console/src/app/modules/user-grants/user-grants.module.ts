import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { InputModule } from '../../modules/input/input.module';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { AvatarModule } from '../avatar/avatar.module';
import { FilterUserGrantsModule } from '../filter-user-grants/filter-user-grants.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { ProjectRoleChipModule } from '../project-role-chip/project-role-chip.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { UserGrantRoleDialogModule } from '../user-grant-role-dialog/user-grant-role-dialog.module';
import { WarnDialogModule } from '../warn-dialog/warn-dialog.module';
import { UserGrantsComponent } from './user-grants.component';

@NgModule({
  declarations: [UserGrantsComponent],
  imports: [
    CommonModule,
    FormsModule,
    AvatarModule,
    MatButtonModule,
    HasRoleModule,
    MatTableModule,
    PaginatorModule,
    MatIconModule,
    RouterModule,
    ProjectRoleChipModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatTooltipModule,
    TableActionsModule,
    UserGrantRoleDialogModule,
    MatSelectModule,
    TranslateModule,
    ActionKeysModule,
    FilterUserGrantsModule,
    HasRolePipeModule,
    TimestampToDatePipeModule,
    RefreshTableModule,
    LocalizedDatePipeModule,
    InputModule,
    WarnDialogModule,
  ],
  exports: [UserGrantsComponent],
})
export class UserGrantsModule {}
