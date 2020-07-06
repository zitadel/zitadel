import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';

import { HttpLoaderFactory } from '../../app.module';
import { HasRoleModule } from '../../directives/has-role/has-role.module';
import { CardModule } from '../../modules/card/card.module';
import { ChangesModule } from '../../modules/changes/changes.module';
import { MetaLayoutModule } from '../../modules/meta-layout/meta-layout.module';
import { ProjectContributorsModule } from '../../modules/project-contributors/project-contributors.module';
import { ProjectRolesModule } from '../../modules/project-roles/project-roles.module';
import { PipesModule } from '../../pipes/pipes.module';
import { GrantedProjectDetailComponent } from './granted-project-detail/granted-project-detail.component';
import { GrantedProjectGridComponent } from './granted-project-grid/granted-project-grid.component';
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
        MatMenuModule,
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
        TranslateModule.forChild({
            loader: {
                provide: TranslateLoader,
                useFactory: HttpLoaderFactory,
                deps: [HttpClient],
            },
        }),
    ],
    schemas: [NO_ERRORS_SCHEMA],
})
export class GrantedProjectsModule { }
