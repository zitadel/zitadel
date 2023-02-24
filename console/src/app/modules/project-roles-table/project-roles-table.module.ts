import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
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
