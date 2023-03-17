import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { TranslateModule } from '@ngx-translate/core';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';

import { FormsModule } from '@angular/forms';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { ToObjectPipeModule } from 'src/app/pipes/to-object/to-object.module';
import { ToPayloadPipeModule } from 'src/app/pipes/to-payload/to-payload.module';
import { DisplayJsonDialogComponent } from './display-json-dialog.component';

@NgModule({
  declarations: [DisplayJsonDialogComponent],
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule,
    MatButtonModule,
    MatIconModule,
    CodemirrorModule,
    TimestampToDatePipeModule,
    ToObjectPipeModule,
    ToPayloadPipeModule,
    LocalizedDatePipeModule,
  ],
  exports: [DisplayJsonDialogComponent],
})
export class DisplayJsonDialogModule {}
