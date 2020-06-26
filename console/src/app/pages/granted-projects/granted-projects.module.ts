import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { NgModule, NgModule } from '@angular/core';
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

import { HttpLoaderFactory } from '../../app.module';
import { HasRoleModule } from '../../directives/has-role/has-role.module';
import { CardModule } from '../../modules/card/card.module';
import { ChangesModule } from '../../modules/changes/changes.module';
import { MetaLayoutModule } from '../../modules/meta-layout/meta-layout.module';
import { ProjectRolesModule } from '../../modules/project-roles/project-roles.module';
import { SearchUserAutocompleteModule } from '../../modules/search-user-autocomplete/search-user-autocomplete.module';
import { PipesModule } from '../../pipes/pipes.module';
import { OrgContributorsModule } from '../orgs/org-contributors/org-contributors.module';
import { UserListModule } from '../user-list/user-list.module';
import { GrantedProjectDetailComponent } from './granted-project-detail/granted-project-detail.component';
import { GrantedProjectGridComponent } from './granted-project-grid/granted-project-grid.component';
import { GrantedProjectListComponent } from './granted-project-list/granted-project-list.component';
import { GrantedProjectsRoutingModule } from './granted-projects-routing.module';

@NgModule({
    declarations: [
        GrantedProjectListComponent,
        GrantedProjectGridComponent,
        GrantedProjectDetailComponent,
    ],
    imports: [
        CommonModule,
        GrantedProjectsRoutingModule,
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
        PipesModule,
        OrgContributorsModule,
        TranslateModule.forChild({
            loader: {
                provide: TranslateLoader,
                useFactory: HttpLoaderFactory,
                deps: [HttpClient],
            },
        }),
    ],
})
export class GrantedProjectsModule { }
