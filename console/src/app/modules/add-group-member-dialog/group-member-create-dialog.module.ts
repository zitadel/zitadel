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

import { GroupMemberCreateDialogComponent } from './group-member-create-dialog.component';

@NgModule({
  declarations: [GroupMemberCreateDialogComponent],
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatChipsModule,
    TranslateModule,
    InputModule,
    MatTooltipModule,
    MatSelectModule,
    FormsModule,
    MatCheckboxModule,
    ReactiveFormsModule,
    SearchUserAutocompleteModule,
  ],
})
export class GroupMemberCreateDialogModule {}
