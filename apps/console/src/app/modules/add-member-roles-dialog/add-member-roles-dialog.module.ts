import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { RoleTransformPipeModule } from 'src/app/pipes/role-transform/role-transform.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { MatDialogModule } from '@angular/material/dialog';
import { AddMemberRolesDialogComponent } from './add-member-roles-dialog.component';

@NgModule({
  declarations: [AddMemberRolesDialogComponent],
  imports: [
    CommonModule,
    TranslateModule,
    MatCheckboxModule,
    MatButtonModule,
    MatDialogModule,
    LocalizedDatePipeModule,
    MatTooltipModule,
    RoleTransformPipeModule,
    TimestampToDatePipeModule,
  ],
})
export class AddMemberRolesDialogModule {}
