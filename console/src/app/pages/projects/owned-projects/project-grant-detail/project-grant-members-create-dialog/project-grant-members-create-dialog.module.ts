import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';

import { ProjectGrantMembersCreateDialogComponent } from './project-grant-members-create-dialog.component';

@NgModule({
    declarations: [ProjectGrantMembersCreateDialogComponent],
    imports: [
        CommonModule,
        FormsModule,
        MatDialogModule,
        MatButtonModule,
        TranslateModule,
        MatSelectModule,
        InputModule,
        SearchUserAutocompleteModule,
    ],
})
export class ProjectGrantMembersCreateDialogModule { }

