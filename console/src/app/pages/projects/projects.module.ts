import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
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
    SharedModule,
    TranslateModule,
    FormsModule,
    HasRoleModule,
    MatTableModule,
    PaginatorModule,
    InputModule,
    MatIconModule,
    MatButtonModule,
    MatProgressBarModule,
    MatCheckboxModule,
    CardModule,
    MatTooltipModule,
    MatSortModule,
    HasRolePipeModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    RefreshTableModule,
    MatRippleModule,
  ],
})
export class ProjectsModule {}
