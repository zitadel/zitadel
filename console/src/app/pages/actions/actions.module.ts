import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ActionTableComponent } from './action-table/action-table.component';
import { ActionsRoutingModule } from './actions-routing.module';
import { ActionsComponent } from './actions.component';
import { AddActionDialogComponent } from './add-action-dialog/add-action-dialog.component';

@NgModule({
  declarations: [
    ActionsComponent,
    ActionTableComponent,
    AddActionDialogComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    ActionsRoutingModule,
    TranslateModule,
    RefreshTableModule,
    MatTableModule,
    PaginatorModule,
    MatButtonModule,
    MatIconModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    MatTooltipModule,
    MatCheckboxModule,
    InputModule,
  ]
})
export class ActionsModule { }
