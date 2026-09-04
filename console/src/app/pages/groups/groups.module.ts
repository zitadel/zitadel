import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { GroupsRoutingModule } from './groups-routing.module';
import { GroupsComponent } from './groups.component';

@NgModule({
  declarations: [GroupsComponent],
  imports: [
    CommonModule,
    GroupsRoutingModule,
    HasRoleModule,
    MatButtonModule,
    MatDialogModule,
    MatIconModule,
    MatTableModule,
    MatTooltipModule,
    PaginatorModule,
    RefreshTableModule,
    TableActionsModule,
    TranslateModule,
    WarnDialogModule,
    LocalizedDatePipeModule,
    TimestampToDatePipeModule,
  ],
})
export default class GroupsModule {}
