import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { CardModule } from '../card/card.module';

import { MatTableModule } from '@angular/material/table';
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
