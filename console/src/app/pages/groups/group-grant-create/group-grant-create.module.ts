import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { CardModule } from 'src/app/modules/card/card.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { SearchProjectAutocompleteModule } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';
import { SearchGroupAutocompleteModule } from 'src/app/modules/search-group-autocomplete/search-group-autocomplete.module';

import { GroupGrantCreateRoutingModule } from './group-grant-create-routing.module';
import { GroupGrantCreateComponent } from './group-grant-create.component';

@NgModule({
  declarations: [GroupGrantCreateComponent],
  imports: [
    GroupGrantCreateRoutingModule,
    CommonModule,
    MatButtonModule,
    MatIconModule,
    CreateLayoutModule,
    TranslateModule,
    CardModule,
    InputModule,
    MatSelectModule,
    SearchProjectAutocompleteModule,
    SearchGroupAutocompleteModule,
    ProjectRolesTableModule,
  ],
})
export default class GroupGrantCreateModule {}
