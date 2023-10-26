import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';

import { FormsModule } from '@angular/forms';
import { MatDialogModule } from '@angular/material/dialog';
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
    MatDialogModule,
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
