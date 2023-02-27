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

import { OverlayModule } from '@angular/cdk/overlay';
import { MatLegacyDialogModule } from '@angular/material/legacy-dialog';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { DisplayJsonDialogModule } from 'src/app/modules/display-json-dialog/display-json-dialog.module';
import { FilterEventsModule } from 'src/app/modules/filter-events/filter-events.module';
import { ToObjectPipeModule } from 'src/app/pipes/to-object/to-object.module';
import { ToPayloadPipeModule } from 'src/app/pipes/to-payload/to-payload.module';
import { EventsRoutingModule } from './events-routing.module';
import { EventsComponent } from './events.component';

@NgModule({
  declarations: [EventsComponent],
  imports: [
    EventsRoutingModule,
    CommonModule,
    TableActionsModule,
    MatIconModule,
    CardModule,
    FilterEventsModule,
    ToObjectPipeModule,
    ToPayloadPipeModule,
    HasRolePipeModule,
    MatLegacyDialogModule,
    MatButtonModule,
    CopyToClipboardModule,
    InputModule,
    TranslateModule,
    InfoSectionModule,
    AvatarModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
    RefreshTableModule,
    ActionKeysModule,
    PaginatorModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    DisplayJsonDialogModule,
    MatTableModule,
    MatSortModule,
    OverlayModule,
  ],
  exports: [],
})
export default class IamViewsModule {}
