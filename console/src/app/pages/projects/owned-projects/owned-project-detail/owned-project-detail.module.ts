import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTabsModule as MatTabsModule } from '@angular/material/legacy-tabs';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { AppCardModule } from 'src/app/modules/app-card/app-card.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { InfoRowModule } from 'src/app/modules/info-row/info-row.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SidenavModule } from 'src/app/modules/sidenav/sidenav.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import ProjectGrantsModule from '../project-grants/project-grants.module';
import ProjectRolesModule from '../project-roles/project-roles.module';
import { ApplicationGridComponent } from './application-grid/application-grid.component';
import { ApplicationsComponent } from './applications/applications.component';
import { OwnedProjectDetailRoutingModule } from './owned-project-detail-routing.module';
import { OwnedProjectDetailComponent } from './owned-project-detail.component';

@NgModule({
  declarations: [OwnedProjectDetailComponent, ApplicationGridComponent, ApplicationsComponent],
  imports: [
    CommonModule,
    FormsModule,
    AppCardModule,
    OwnedProjectDetailRoutingModule,
    TranslateModule,
    ReactiveFormsModule,
    HasRoleModule,
    MatButtonModule,
    MatIconModule,
    InfoRowModule,
    ContributorsModule,
    MatTabsModule,
    WarnDialogModule,
    MatTooltipModule,
    ProjectRolesTableModule,
    HasRolePipeModule,
    UserGrantsModule,
    TimestampToDatePipeModule,
    SidenavModule,
    MatTableModule,
    InputModule,
    CardModule,
    PaginatorModule,
    ProjectGrantsModule,
    ProjectRolesModule,
    MatRippleModule,
    TopViewModule,
    MatCheckboxModule,
    MatSelectModule,
    InfoSectionModule,
    MatMenuModule,
    MatProgressSpinnerModule,
    ChangesModule,
    MetaLayoutModule,
    RefreshTableModule,
    MemberCreateDialogModule,
    LocalizedDatePipeModule,
  ],
})
export default class OwnedProjectDetailModule {}
