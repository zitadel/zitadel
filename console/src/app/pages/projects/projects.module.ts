import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { CUSTOM_ELEMENTS_SCHEMA, NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialogModule } from '@angular/material/dialog';
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
import { HttpLoaderFactory } from 'src/app/app.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';

import { ChangesModule } from '../../modules/changes/changes.module';
import { ProjectRolesModule } from '../../modules/project-roles/project-roles.module';
import { OrgMembersModule } from '../orgs/org-members/org-members.module';
import { UserListModule } from '../user-list/user-list.module';
import { ProjectApplicationGridComponent } from './project-application-grid/project-application-grid.component';
import { ProjectApplicationsComponent } from './project-applications/project-applications.component';
import { ProjectContributorsComponent } from './project-contributors/project-contributors.component';
import { ProjectDetailComponent } from './project-detail/project-detail.component';
import {
    ProjectGrantMembersCreateDialogComponent,
} from './project-grant-members-create-dialog/project-grant-members-create-dialog.component';
import { ProjectGrantMembersComponent } from './project-grant-members/project-grant-members.component';
import { ProjectGrantsComponent } from './project-grants/project-grants.component';
import { ProjectGridComponent } from './project-grid/project-grid.component';
import { ProjectListComponent } from './project-list/project-list.component';
import { ProjectsRoutingModule } from './projects-routing.module';

@NgModule({
    declarations: [
        ProjectListComponent,
        ProjectDetailComponent,
        ProjectApplicationsComponent,
        ProjectGridComponent,
        ProjectApplicationGridComponent,
        ProjectGrantsComponent,
        ProjectGrantMembersComponent,
        ProjectGrantMembersCreateDialogComponent,
        ProjectContributorsComponent,
    ],
    imports: [
        ProjectsRoutingModule,
        CommonModule,
        FormsModule,
        TranslateModule,
        ReactiveFormsModule,
        HasRoleModule,
        MatTableModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatInputModule,
        ChangesModule,
        UserListModule,
        MatMenuModule,
        MatChipsModule,
        MatIconModule,
        MatSelectModule,
        MatButtonModule,
        MatProgressSpinnerModule,
        MetaLayoutModule,
        MatProgressBarModule,
        MatDialogModule,
        MatButtonToggleModule,
        MatTabsModule,
        ProjectRolesModule,
        SearchUserAutocompleteModule,
        MatCheckboxModule,
        CardModule,
        MatTooltipModule,
        MatSortModule,
        OrgMembersModule,
        TranslateModule.forChild({
            loader: {
                provide: TranslateLoader,
                useFactory: HttpLoaderFactory,
                deps: [HttpClient],
            },
        }),
    ],
    entryComponents: [
        ProjectGrantMembersCreateDialogComponent,
    ],
    exports: [ProjectListComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA],
})
export class ProjectsModule { }
