import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { TypeSafeCellDefModule } from 'src/app/directives/type-safe-cell-def/type-safe-cell-def.module';

import { AddTokenDialogModule } from '../add-token-dialog/add-token-dialog.module';
import { CardModule } from '../card/card.module';
import { InputModule } from '../input/input.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { ShowTokenDialogModule } from '../show-token-dialog/show-token-dialog.module';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { WarnDialogModule } from '../warn-dialog/warn-dialog.module';
import { PersonalAccessTokensComponent } from './personal-access-tokens.component';

@NgModule({
  declarations: [PersonalAccessTokensComponent],
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    MatButtonModule,
    MatDialogModule,
    HasRoleModule,
    CardModule,
    MatTableModule,
    PaginatorModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatTooltipModule,
    HasRolePipeModule,
    TimestampToDatePipeModule,
    TableActionsModule,
    LocalizedDatePipeModule,
    TranslateModule,
    RefreshTableModule,
    InputModule,
    ShowTokenDialogModule,
    WarnDialogModule,
    AddTokenDialogModule,
    TypeSafeCellDefModule,
  ],
  exports: [PersonalAccessTokensComponent],
})
export class PersonalAccessTokensModule {}
