import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
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
import { FilterGroupGrantsModule } from '../filter-group-grants/filter-group-grants.module';
import { InputModule } from '../../modules/input/input.module';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { AvatarModule } from '../avatar/avatar.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { ProjectRoleChipModule } from '../project-role-chip/project-role-chip.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { UserGrantRoleDialogModule } from '../user-grant-role-dialog/user-grant-role-dialog.module';
import { WarnDialogModule } from '../warn-dialog/warn-dialog.module';
import { GroupGrantsComponent } from './group-grants.component';

@NgModule({
  declarations: [GroupGrantsComponent],
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
    HasRolePipeModule,
    TimestampToDatePipeModule,
    FilterGroupGrantsModule,
    RefreshTableModule,
    LocalizedDatePipeModule,
    InputModule,
    WarnDialogModule,
  ],
  exports: [GroupGrantsComponent],
})
export class GroupGrantsModule {}
