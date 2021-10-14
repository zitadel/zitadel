import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { DurationToSecondsPipeModule } from 'src/app/pipes/duration-to-seconds-pipe/duration-to-seconds-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ActionTableComponent } from './action-table/action-table.component';
import { ActionsRoutingModule } from './actions-routing.module';
import { ActionsComponent } from './actions.component';
import { AddActionDialogComponent } from './add-action-dialog/add-action-dialog.component';
import { AddFlowDialogComponent } from './add-flow-dialog/add-flow-dialog.component';

@NgModule({
  declarations: [
    ActionsComponent,
    ActionTableComponent,
    AddActionDialogComponent,
    AddFlowDialogComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    ActionsRoutingModule,
    TranslateModule,
    MatDialogModule,
    RefreshTableModule,
    MatTableModule,
    PaginatorModule,
    MatButtonModule,
    ReactiveFormsModule,
    MatIconModule,
    DurationToSecondsPipeModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    MatTooltipModule,
    MatCheckboxModule,
    InputModule,
    FormFieldModule,
    MatSelectModule,
  ]
})
export class ActionsModule { }
