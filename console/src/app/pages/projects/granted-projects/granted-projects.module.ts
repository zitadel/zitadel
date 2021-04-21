import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { ContributorsModule } from 'src/app/modules/contributors/contributors.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { ProjectRolesModule } from 'src/app/modules/project-roles/project-roles.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { GrantedProjectDetailComponent } from './granted-project-detail/granted-project-detail.component';
import { GrantedProjectGridComponent } from './granted-project-list/granted-project-grid/granted-project-grid.component';
import { GrantedProjectListComponent } from './granted-project-list/granted-project-list.component';
import { GrantedProjectsRoutingModule } from './granted-projects-routing.module';
import { GrantedProjectsComponent } from './granted-projects.component';

@NgModule({
    declarations: [
        GrantedProjectsComponent,
        GrantedProjectListComponent,
        GrantedProjectGridComponent,
        GrantedProjectDetailComponent,
    ],
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
        ProjectRolesModule,
        MatCheckboxModule,
        CardModule,
        MatTooltipModule,
        MatSortModule,
        HasRolePipeModule,
        TranslateModule,
        TimestampToDatePipeModule,
        SharedModule,
        LocalizedDatePipeModule,
        MemberCreateDialogModule,
        RefreshTableModule,
    ],
})
export class GrantedProjectsModule { }
