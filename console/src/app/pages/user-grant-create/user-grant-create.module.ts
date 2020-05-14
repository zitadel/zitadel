import { CommonModule } from '@angular/common';
import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';
import { CardModule } from 'src/app/modules/card/card.module';
import {
    SearchProjectAutocompleteModule,
} from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';

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
        SearchProjectAutocompleteModule,
        ProjectRolesModule,
    ],
    schemas: [
        CUSTOM_ELEMENTS_SCHEMA,
    ],
})
export class UserGrantCreateModule { }
