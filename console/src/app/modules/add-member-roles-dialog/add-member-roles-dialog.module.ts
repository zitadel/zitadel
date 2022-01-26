import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { AddMemberRolesDialogComponent } from './add-member-roles-dialog.component';

@NgModule({
  declarations: [AddMemberRolesDialogComponent],
  imports: [CommonModule, TranslateModule, MatButtonModule, LocalizedDatePipeModule, TimestampToDatePipeModule],
})
export class AddMemberRolesDialogModule {}
