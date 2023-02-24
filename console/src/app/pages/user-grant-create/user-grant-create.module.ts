import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { TranslateModule } from '@ngx-translate/core';
import { CardModule } from 'src/app/modules/card/card.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { SearchProjectAutocompleteModule } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';

import { UserGrantCreateRoutingModule } from './user-grant-create-routing.module';
import { UserGrantCreateComponent } from './user-grant-create.component';

@NgModule({
  declarations: [UserGrantCreateComponent],
  imports: [
    UserGrantCreateRoutingModule,
    CommonModule,
    MatButtonModule,
    MatIconModule,
    CreateLayoutModule,
    TranslateModule,
    CardModule,
    InputModule,
    MatSelectModule,
    SearchProjectAutocompleteModule,
    SearchUserAutocompleteModule,
    ProjectRolesTableModule,
  ],
})
export default class UserGrantCreateModule {}
