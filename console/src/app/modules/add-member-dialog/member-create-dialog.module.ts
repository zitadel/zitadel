import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';

import {
    OrgMemberRolesAutocompleteModule,
} from '../../pages/orgs/org-member-roles-autocomplete/org-member-roles-autocomplete.module';
import { SearchProjectAutocompleteModule } from '../search-project-autocomplete/search-project-autocomplete.module';
import { SearchRolesAutocompleteModule } from '../search-roles-autocomplete/search-roles-autocomplete.module';
import { MemberCreateDialogComponent } from './member-create-dialog.component';

@NgModule({
    declarations: [MemberCreateDialogComponent],
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        MatChipsModule,
        MatInputModule,
        TranslateModule,
        MatFormFieldModule,
        MatSelectModule,
        FormsModule,
        ReactiveFormsModule,
        SearchUserAutocompleteModule,
        SearchRolesAutocompleteModule,
        SearchProjectAutocompleteModule,
        OrgMemberRolesAutocompleteModule,
    ],
})
export class MemberCreateDialogModule { }
