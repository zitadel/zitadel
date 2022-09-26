import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { InfoRowComponent } from './info-row.component';

@NgModule({
  declarations: [InfoRowComponent],
  imports: [
    CommonModule,
    FormsModule,
    MatTooltipModule,
    TranslateModule,
    CopyToClipboardModule,
    MatButtonModule,
    LocalizedDatePipeModule,
    TimestampToDatePipeModule,
  ],
  exports: [InfoRowComponent],
})
export class InfoRowModule {}
