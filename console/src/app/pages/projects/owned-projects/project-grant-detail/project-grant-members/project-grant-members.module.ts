import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';

import {
    ProjectGrantMembersCreateDialogModule,
} from './project-grant-members-create-dialog/project-grant-members-create-dialog.module';
import { ProjectGrantMembersComponent } from './project-grant-members.component';

@NgModule({
    declarations: [ProjectGrantMembersComponent],
    imports: [
        CommonModule,
        HasRoleModule,
        RouterModule,
        MatButtonModule,
        MatCheckboxModule,
        MatIconModule,
        MatInputModule,
        MatFormFieldModule,
        MatSelectModule,
        MatTableModule,
        SearchUserAutocompleteModule,
        ProjectGrantMembersCreateDialogModule,
        MatPaginatorModule,
        MatSortModule,
        MatTooltipModule,
        MatDialogModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
        RefreshTableModule,
        HasRolePipeModule,
    ],
    exports: [
        ProjectGrantMembersComponent,
    ],
})
export class ProjectGrantMembersModule { }
