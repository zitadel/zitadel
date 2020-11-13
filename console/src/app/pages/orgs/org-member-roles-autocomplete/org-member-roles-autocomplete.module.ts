import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

import { OrgMemberRolesAutocompleteComponent } from './org-member-roles-autocomplete.component';

@NgModule({
    declarations: [OrgMemberRolesAutocompleteComponent],
    imports: [
        CommonModule,
        MatButtonModule,
        MatSelectModule,
        FormFieldModule,
        MatIconModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
    ],
    exports: [
        OrgMemberRolesAutocompleteComponent,
    ],
})
export class OrgMemberRolesAutocompleteModule { }
