import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { TruncatePipeModule } from 'src/app/pipes/truncate-pipe/truncate-pipe.module';

import { PaginatorModule } from '../paginator/paginator.module';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { SMTPTableComponent } from './smtp-table.component';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatTableModule } from '@angular/material/table';
import { SmtpTestDialogModule } from '../smtp-test-dialog/smtp-test-dialog.module';

@NgModule({
  declarations: [SMTPTableComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MatButtonModule,
    TableActionsModule,
    MatCheckboxModule,
    MatIconModule,
    MatTooltipModule,
    TranslateModule,
    LocalizedDatePipeModule,
    TimestampToDatePipeModule,
    MatTableModule,
    PaginatorModule,
    RouterModule,
    RefreshTableModule,
    HasRoleModule,
    HasRolePipeModule,
    TruncatePipeModule,
    SmtpTestDialogModule,
  ],
  exports: [SMTPTableComponent],
})
export class SMTPTableModule {}
