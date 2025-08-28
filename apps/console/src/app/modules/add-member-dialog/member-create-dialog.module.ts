import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';
import { RoleTransformPipeModule } from 'src/app/pipes/role-transform/role-transform.module';

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
    TranslateModule,
    InputModule,
    MatTooltipModule,
    MatSelectModule,
    RoleTransformPipeModule,
    FormsModule,
    MatCheckboxModule,
    ReactiveFormsModule,
    SearchUserAutocompleteModule,
    SearchRolesAutocompleteModule,
    SearchProjectAutocompleteModule,
  ],
})
export class MemberCreateDialogModule {}
