import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActionsTwoActionsComponent } from './actions-two-actions/actions-two-actions.component';
import { ActionsTwoTargetsComponent } from './actions-two-targets/actions-two-targets.component';
import { ActionsTwoRoutingModule } from './actions-two-routing.module';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { ActionsTwoTargetsTableComponent } from './actions-two-targets/actions-two-targets-table/actions-two-targets-table.component';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { ActionsTwoActionsTableComponent } from './actions-two-actions/actions-two-actions-table/actions-two-actions-table.component';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { TypeSafeCellDefModule } from 'src/app/directives/type-safe-cell-def/type-safe-cell-def.module';
import { ProjectRoleChipModule } from '../project-role-chip/project-role-chip.module';
import { ActionConditionPipeModule } from 'src/app/pipes/action-condition-pipe/action-condition-pipe.module';
import { MatSelectModule } from '@angular/material/select';
import { MatIconModule } from '@angular/material/icon';
import { InfoSectionModule } from '../info-section/info-section.module';

@NgModule({
  declarations: [
    ActionsTwoActionsComponent,
    ActionsTwoTargetsComponent,
    ActionsTwoTargetsTableComponent,
    ActionsTwoActionsTableComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    TableActionsModule,
    TimestampToDatePipeModule,
    ActionsTwoRoutingModule,
    LocalizedDatePipeModule,
    ReactiveFormsModule,
    TranslateModule,
    MatTableModule,
    MatTooltipModule,
    MatSelectModule,
    RefreshTableModule,
    ActionKeysModule,
    MatIconModule,
    TypeSafeCellDefModule,
    ProjectRoleChipModule,
    ActionConditionPipeModule,
    InfoSectionModule,
  ],
  exports: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent, ActionsTwoTargetsTableComponent],
})
export default class ActionsTwoModule {}
