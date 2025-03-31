import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
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
import { WarnDialogModule } from '../warn-dialog/warn-dialog.module';
import { GroupCreateDialogModule } from '../add-group-dialog/group-create-dialog.module';
import { UserGroupsComponent } from './user-groups.component';

@NgModule({
  declarations: [UserGroupsComponent],
  imports: [
    CommonModule,
    FormsModule,
    AvatarModule,
    MatButtonModule,
    HasRoleModule,
    MatTableModule,
    MatDialogModule,
    PaginatorModule,
    MatIconModule,
    RouterModule,
    ProjectRoleChipModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatTooltipModule,
    TableActionsModule,
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
    GroupCreateDialogModule,
  ],
  exports: [UserGroupsComponent],
})
export class UserGroupsModule {}
