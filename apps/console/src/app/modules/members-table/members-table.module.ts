import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
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
