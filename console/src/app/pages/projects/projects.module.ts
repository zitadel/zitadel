import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { FilterProjectModule } from 'src/app/modules/filter-project/filter-project.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { NavToggleModule } from 'src/app/modules/nav-toggle/nav-toggle.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ProjectGridComponent } from './project-grid/project-grid.component';
import { ProjectListComponent } from './project-list/project-list.component';
import { ProjectsRoutingModule } from './projects-routing.module';
import { ProjectsComponent } from './projects.component';

@NgModule({
  declarations: [ProjectsComponent, ProjectListComponent, ProjectGridComponent],
  imports: [
    CommonModule,
    ProjectsRoutingModule,
    TranslateModule,
    FormsModule,
    HasRoleModule,
    MatTableModule,
    PaginatorModule,
    InputModule,
    MatIconModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    CardModule,
    MatTooltipModule,
    FilterProjectModule,
    ActionKeysModule,
    TableActionsModule,
    MatSortModule,
    HasRolePipeModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    RefreshTableModule,
    MatRippleModule,
    NavToggleModule,
  ],
})
export default class ProjectsModule {}
