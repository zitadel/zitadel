import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { ActionKeysModule } from '../action-keys/action-keys.module';
import { RefreshTableComponent } from './refresh-table.component';

@NgModule({
  declarations: [RefreshTableComponent],
  imports: [
    CommonModule,
    MatButtonModule,
    MatIconModule,
    TranslateModule,
    ActionKeysModule,
    FormsModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    PaginatorModule,
  ],
  exports: [RefreshTableComponent],
})
export class RefreshTableModule {}
