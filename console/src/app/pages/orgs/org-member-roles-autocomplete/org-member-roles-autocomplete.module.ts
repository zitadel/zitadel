import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { OrgMemberRolesAutocompleteComponent } from './org-member-roles-autocomplete.component';

@NgModule({
    declarations: [OrgMemberRolesAutocompleteComponent],
    imports: [
        CommonModule,
        MatButtonModule,
        MatSelectModule,
        InputModule,
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
