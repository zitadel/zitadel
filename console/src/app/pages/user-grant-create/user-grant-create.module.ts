import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { CardModule } from 'src/app/modules/card/card.module';
import { InputModule } from 'src/app/modules/input/input.module';
import {
    SearchProjectAutocompleteModule,
} from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';

import { ProjectRolesModule } from '../../modules/project-roles/project-roles.module';
import { UserGrantCreateRoutingModule } from './user-grant-create-routing.module';
import { UserGrantCreateComponent } from './user-grant-create.component';



@NgModule({
    declarations: [UserGrantCreateComponent],
    imports: [
        UserGrantCreateRoutingModule,
        CommonModule,
        MatButtonModule,
        MatIconModule,
        TranslateModule,
        CardModule,
        InputModule,
        MatSelectModule,
        SearchProjectAutocompleteModule,
        SearchUserAutocompleteModule,
        ProjectRolesModule,
        SharedModule,
    ],
})
export class UserGrantCreateModule { }
