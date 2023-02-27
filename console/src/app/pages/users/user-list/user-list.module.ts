import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { FilterUserModule } from 'src/app/modules/filter-user/filter-user.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { NavToggleModule } from 'src/app/modules/nav-toggle/nav-toggle.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { UserListComponent } from './user-list.component';
import { UserTableComponent } from './user-table/user-table.component';

@NgModule({
  declarations: [UserListComponent, UserTableComponent],
  imports: [
    AvatarModule,
    CommonModule,
    FormsModule,
    MatButtonModule,
    MatDialogModule,
    HasRoleModule,
    CardModule,
    MatTableModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatTooltipModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    HasRolePipeModule,
    TranslateModule,
    FilterUserModule,
    RouterModule,
    NavToggleModule,
    RefreshTableModule,
    TableActionsModule,
    ActionKeysModule,
    MatMenuModule,
    MatSortModule,
    InputModule,
    PaginatorModule,
  ],
  exports: [UserListComponent],
})
export class UserListModule {}
