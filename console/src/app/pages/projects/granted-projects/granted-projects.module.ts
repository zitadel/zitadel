import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTabsModule as MatTabsModule } from '@angular/material/legacy-tabs';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatSortModule } from '@angular/material/sort';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { InfoRowModule } from 'src/app/modules/info-row/info-row.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { GrantedProjectDetailComponent } from './granted-project-detail/granted-project-detail.component';
import { GrantedProjectsRoutingModule } from './granted-projects-routing.module';

@NgModule({
  declarations: [GrantedProjectDetailComponent],
  imports: [
    CommonModule,
    UserGrantsModule,
    GrantedProjectsRoutingModule,
    ContributorsModule,
    FormsModule,
    TranslateModule,
    ReactiveFormsModule,
    HasRoleModule,
    MatTableModule,
    PaginatorModule,
    InputModule,
    ChangesModule,
    MatIconModule,
    MatSelectModule,
    MatButtonModule,
    MatTabsModule,
    MatProgressSpinnerModule,
    MetaLayoutModule,
    MatProgressBarModule,
    ProjectRolesTableModule,
    MatCheckboxModule,
    CardModule,
    MatTooltipModule,
    MatSortModule,
    HasRolePipeModule,
    TimestampToDatePipeModule,
    TopViewModule,
    InfoRowModule,
    LocalizedDatePipeModule,
    MemberCreateDialogModule,
    MatRippleModule,
    RefreshTableModule,
  ],
})
export default class GrantedProjectsModule {}
