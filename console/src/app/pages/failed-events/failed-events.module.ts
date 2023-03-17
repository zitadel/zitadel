import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { FailedEventsRoutingModule } from './failed-events-routing.module';
import { FailedEventsComponent } from './failed-events.component';

@NgModule({
  declarations: [FailedEventsComponent],
  imports: [
    FailedEventsRoutingModule,
    CommonModule,
    TableActionsModule,
    MatIconModule,
    CardModule,
    HasRolePipeModule,
    PaginatorModule,
    MatButtonModule,
    CopyToClipboardModule,
    InputModule,
    TranslateModule,
    InfoSectionModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    RefreshTableModule,
    LocalizedDatePipeModule,
    TimestampToDatePipeModule,
    MatTableModule,
    MatSortModule,
  ],
})
export default class FailedEventsModule {}
