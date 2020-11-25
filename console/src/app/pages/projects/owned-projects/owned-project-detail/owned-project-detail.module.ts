import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
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
import { ProjectRolesModule } from 'src/app/modules/project-roles/project-roles.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

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
        FormsModule,
        OwnedProjectDetailRoutingModule,
        TranslateModule,
        ReactiveFormsModule,
        HasRoleModule,
        MatButtonModule,
        MatIconModule,
        ContributorsModule,
        MatTabsModule,
        WarnDialogModule,
        MatTooltipModule,
        ProjectRolesModule,
        HasRolePipeModule,
        UserGrantsModule,
        TimestampToDatePipeModule,
        MatTableModule,
        InputModule,
        CardModule,
        MatPaginatorModule,
        MatRippleModule,
        MatCheckboxModule,
        MatSelectModule,
        MatProgressSpinnerModule,
        ChangesModule,
        MetaLayoutModule,
        RefreshTableModule,
        MemberCreateDialogModule,
        LocalizedDatePipeModule,
    ],
})
export class OwnedProjectDetailModule { }
