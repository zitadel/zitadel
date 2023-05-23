import { DragDropModule } from '@angular/cdk/drag-drop';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { DurationToSecondsPipeModule } from 'src/app/pipes/duration-to-seconds-pipe/duration-to-seconds-pipe.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ActionTableComponent } from './action-table/action-table.component';
import { ActionsRoutingModule } from './actions-routing.module';
import { ActionsComponent } from './actions.component';
import { AddActionDialogComponent } from './add-action-dialog/add-action-dialog.component';
import { AddFlowDialogComponent } from './add-flow-dialog/add-flow-dialog.component';

@NgModule({
  declarations: [ActionsComponent, ActionTableComponent, AddActionDialogComponent, AddFlowDialogComponent],
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
    HasRoleModule,
    ActionKeysModule,
    MatTooltipModule,
    CardModule,
    MatCheckboxModule,
    InputModule,
    FormFieldModule,
    MatSelectModule,
    WarnDialogModule,
    DragDropModule,
    InfoSectionModule,
    HasRolePipeModule,
    TableActionsModule,
    CodemirrorModule,
  ],
})
export default class ActionsModule {}
