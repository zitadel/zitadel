import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { ProjectRoleChipModule } from 'src/app/modules/project-role-chip/project-role-chip.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ProjectGrantsComponent } from './project-grants.component';

@NgModule({
  declarations: [ProjectGrantsComponent],
  imports: [
    CommonModule,
    FormsModule,
    TimestampToDatePipeModule,
    TableActionsModule,
    MatTooltipModule,
    MatCheckboxModule,
    RefreshTableModule,
    RouterModule,
    LocalizedDatePipeModule,
    ProjectRoleChipModule,
    MatButtonModule,
    HasRolePipeModule,
    MatIconModule,
    InputModule,
    MatTableModule,
    TranslateModule,
    ActionKeysModule,
    MatSelectModule,
    PaginatorModule,
  ],
  exports: [ProjectGrantsComponent],
})
export default class ProjectGrantsModule {}
