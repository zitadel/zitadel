import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';
import { RoleTransformPipeModule } from 'src/app/pipes/role-transform/role-transform.module';

import { AddMemberRolesDialogModule } from '../add-member-roles-dialog/add-member-roles-dialog.module';
import { AvatarModule } from '../avatar/avatar.module';
import { PaginatorModule } from '../paginator/paginator.module';
import { RefreshTableModule } from '../refresh-table/refresh-table.module';
import { TableActionsModule } from '../table-actions/table-actions.module';
import { WarnDialogModule } from '../warn-dialog/warn-dialog.module';
import { MembersTableComponent } from './members-table.component';

@NgModule({
  declarations: [MembersTableComponent],
  imports: [
    CommonModule,
    InputModule,
    MatSelectModule,
    MatCheckboxModule,
    MatIconModule,
    MatTableModule,
    MatChipsModule,
    RoleTransformPipeModule,
    PaginatorModule,
    AddMemberRolesDialogModule,
    MatSortModule,
    MatTooltipModule,
    FormsModule,
    TranslateModule,
    WarnDialogModule,
    RefreshTableModule,
    TableActionsModule,
    RouterModule,
    AvatarModule,
    MatButtonModule,
  ],
  exports: [MembersTableComponent],
})
export class MembersTableModule {}
