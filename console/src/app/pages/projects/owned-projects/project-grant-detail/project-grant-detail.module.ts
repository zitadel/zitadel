import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyAutocompleteModule as MatAutocompleteModule } from '@angular/material/legacy-autocomplete';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MembersTableModule } from 'src/app/modules/members-table/members-table.module';
import { ProjectRoleChipModule } from 'src/app/modules/project-role-chip/project-role-chip.module';
import { UserGrantRoleDialogModule } from 'src/app/modules/user-grant-role-dialog/user-grant-role-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { ProjectGrantDetailRoutingModule } from './project-grant-detail-routing.module';
import { ProjectGrantDetailComponent } from './project-grant-detail.component';
import { ProjectGrantIllustrationComponent } from './project-grant-illustration/project-grant-illustration.component';

@NgModule({
  declarations: [ProjectGrantDetailComponent, ProjectGrantIllustrationComponent],
  imports: [
    CommonModule,
    ProjectGrantDetailRoutingModule,
    MatAutocompleteModule,
    HasRoleModule,
    MatChipsModule,
    MatButtonModule,
    MatCheckboxModule,
    MatMenuModule,
    UserGrantRoleDialogModule,
    MatIconModule,
    MatTableModule,
    InputModule,
    MatTooltipModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    ActionKeysModule,
    ProjectRoleChipModule,
    FormsModule,
    TranslateModule,
    MatSelectModule,
    DetailLayoutModule,
    MemberCreateDialogModule,
    HasRolePipeModule,
    MembersTableModule,
    MatDialogModule,
  ],
})
export default class ProjectGrantDetailModule {}
