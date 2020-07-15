import { CommonModule } from '@angular/common';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { ProjectContributorsModule } from 'src/app/modules/project-contributors/project-contributors.module';
import { ProjectRolesModule } from 'src/app/modules/project-roles/project-roles.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe.module';

import { ProjectGrantMembersModule } from '../project-grant-detail/project-grant-members/project-grant-members.module';
import { ApplicationGridComponent } from './application-grid/application-grid.component';
import { ApplicationsComponent } from './applications/applications.component';
import { OwnedProjectDetailRoutingModule } from './owned-project-detail-routing.module';
import { OwnedProjectDetailComponent } from './owned-project-detail.component';
import { ProjectGrantsComponent } from './project-grants/project-grants.component';

@NgModule({
    declarations: [
        OwnedProjectDetailComponent,
        ApplicationGridComponent,
        ApplicationsComponent,
        ProjectGrantsComponent,
    ],
    imports: [
        CommonModule,
        OwnedProjectDetailRoutingModule,
        TranslateModule,
        HasRoleModule,
        MatTabsModule,
        MatButtonModule,
        MatIconModule,
        MetaLayoutModule,
        ProjectContributorsModule,
        WarnDialogModule,
        ProjectRolesModule,
        HasRolePipeModule,
        ProjectGrantMembersModule,
        TimestampToDatePipeModule,
    ],
    schemas: [NO_ERRORS_SCHEMA],
})
export class OwnedProjectDetailModule { }
