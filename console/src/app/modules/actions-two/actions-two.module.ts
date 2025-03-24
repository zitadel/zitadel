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
import { TypeSafeCellDefDirective } from './actions-two-targets/actions-two-targets-table/type-safe-cell-def.directive';
import { ActionsTwoActionsTableComponent } from './actions-two-actions/actions-two-actions-table/actions-two-actions-table.component';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

@NgModule({
  declarations: [
    TypeSafeCellDefDirective,
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
    RefreshTableModule,
    ActionKeysModule,
  ],
  exports: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent, ActionsTwoTargetsTableComponent],
})
export default class ActionsTwoModule {}
