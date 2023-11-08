import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
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
