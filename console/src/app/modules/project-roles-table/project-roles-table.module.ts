import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ActionKeysModule } from '../action-keys/action-keys.module';
import { ProjectRoleChipModule } from '../project-role-chip/project-role-chip.module';
import { ProjectRoleDetailDialogModule } from '../project-role-detail-dialog/project-role-detail-dialog.module';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { ProjectRolesTableComponent } from './project-roles-table.component';

@NgModule({
  declarations: [ProjectRolesTableComponent],
  imports: [
    CommonModule,
    MatButtonModule,
    ProjectRoleChipModule,
    HasRoleModule,
    MatTableModule,
    PaginatorModule,
    MatDialogModule,
    InputModule,
    RouterModule,
    FormsModule,
    ActionKeysModule,
    ReactiveFormsModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatTooltipModule,
    HasRolePipeModule,
    TranslateModule,
    TableActionsModule,
    ProjectRoleDetailDialogModule,
    MatMenuModule,
    TimestampToDatePipeModule,
    RefreshTableModule,
    LocalizedDatePipeModule,
  ],
  exports: [ProjectRolesTableComponent],
})
export class ProjectRolesTableModule {}
