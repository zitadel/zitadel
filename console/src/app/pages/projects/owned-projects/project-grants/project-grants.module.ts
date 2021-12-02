import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ProjectGrantsRoutingModule } from './project-grants-routing.module';
import { ProjectGrantsComponent } from './project-grants.component';

@NgModule({
  declarations: [ProjectGrantsComponent],
  imports: [
    CommonModule,
    FormsModule,
    ProjectGrantsRoutingModule,
    TimestampToDatePipeModule,
    MatCheckboxModule,
    RefreshTableModule,
    LocalizedDatePipeModule,
    MatButtonModule,
    HasRolePipeModule,
    MatIconModule,
    InputModule,
    MatTableModule,
    TranslateModule,
    MatSelectModule,
    PaginatorModule,
  ],
})
export class ProjectGrantsModule {}
