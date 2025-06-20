import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
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
import { NavToggleModule } from 'src/app/modules/nav-toggle/nav-toggle.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SidenavModule } from 'src/app/modules/sidenav/sidenav.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { GroupGrantsModule } from 'src/app/modules/group-grants/group-grants.module';
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
    WarnDialogModule,
    MatTooltipModule,
    ProjectRolesTableModule,
    HasRolePipeModule,
    UserGrantsModule,
    GroupGrantsModule,
    TimestampToDatePipeModule,
    SidenavModule,
    MatTableModule,
    InputModule,
    CardModule,
    PaginatorModule,
    ProjectGrantsModule,
    ProjectRolesModule,
    NavToggleModule,
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
