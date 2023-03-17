import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { MembersTableModule } from 'src/app/modules/members-table/members-table.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { OrgMembersRoutingModule } from './org-members-routing.module';
import { OrgMembersComponent } from './org-members.component';

@NgModule({
  declarations: [OrgMembersComponent],
  imports: [
    OrgMembersRoutingModule,
    CommonModule,
    MatChipsModule,
    MatButtonModule,
    ActionKeysModule,
    HasRoleModule,
    MatIconModule,
    MatTooltipModule,
    TranslateModule,
    DetailLayoutModule,
    RefreshTableModule,
    MembersTableModule,
    HasRolePipeModule,
    MemberCreateDialogModule,
  ],
})
export default class OrgMembersModule {}
