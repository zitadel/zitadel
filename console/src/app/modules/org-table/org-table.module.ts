import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyRadioModule as MatRadioModule } from '@angular/material/legacy-radio';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ActionKeysModule } from '../action-keys/action-keys.module';
import { FilterOrgModule } from '../filter-org/filter-org.module';
import { InputModule } from '../input/input.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { OrgTableComponent } from './org-table.component';

@NgModule({
  declarations: [OrgTableComponent],
  imports: [
    CommonModule,
    MatTableModule,
    TranslateModule,
    RefreshTableModule,
    ActionKeysModule,
    FilterOrgModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    MatSortModule,
    TableActionsModule,
    MatIconModule,
    PaginatorModule,
    HasRoleModule,
    RouterModule,
    MatButtonModule,
    MatTooltipModule,
    CopyToClipboardModule,
    MatRadioModule,
    InputModule,
    FormsModule,
  ],
  exports: [OrgTableComponent],
})
export class OrgTableModule {}
