import { CommonModule } from '@angular/common';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { ProjectContributorsModule } from 'src/app/modules/project-contributors/project-contributors.module';
import { ProjectRolesModule } from 'src/app/modules/project-roles/project-roles.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';
import { PipesModule } from 'src/app/pipes/pipes.module';

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
        ProjectContributorsModule,
        FormsModule,
        TranslateModule,
        ReactiveFormsModule,
        HasRoleModule,
        MatTableModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatInputModule,
        ChangesModule,
        MatIconModule,
        MatSelectModule,
        MatButtonModule,
        MatProgressSpinnerModule,
        MetaLayoutModule,
        MatProgressBarModule,
        MatTabsModule,
        ProjectRolesModule,
        MatCheckboxModule,
        CardModule,
        MatTooltipModule,
        MatSortModule,
        PipesModule,
        HasRolePipeModule,
        TranslateModule,
    ],
    schemas: [NO_ERRORS_SCHEMA],
})
export class GrantedProjectsModule { }
