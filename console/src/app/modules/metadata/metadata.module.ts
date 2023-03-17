import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { CardModule } from '../card/card.module';

import { InputModule } from '../input/input.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { MetadataDialogComponent } from './metadata-dialog/metadata-dialog.component';
import { MetadataComponent } from './metadata/metadata.component';

@NgModule({
  declarations: [MetadataComponent, MetadataDialogComponent],
  imports: [
    CommonModule,
    MatDialogModule,
    MatProgressSpinnerModule,
    CardModule,
    MatButtonModule,
    TranslateModule,
    InputModule,
    MatIconModule,
    MatTooltipModule,
    FormsModule,
    LocalizedDatePipeModule,
    TimestampToDatePipeModule,
    RefreshTableModule,
    MatTableModule,
  ],
  exports: [MetadataComponent, MetadataDialogComponent],
})
export class MetadataModule {}
