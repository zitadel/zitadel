import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CardModule } from 'src/app/modules/card/card.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { SearchOrgAutocompleteModule } from 'src/app/modules/search-org-autocomplete/search-org-autocomplete.module';
import { ProjectGrantCreateRoutingModule } from './project-grant-create-routing.module';
import { ProjectGrantCreateComponent } from './project-grant-create.component';

@NgModule({
  declarations: [ProjectGrantCreateComponent],
  imports: [
    ProjectGrantCreateRoutingModule,
    CommonModule,
    MatAutocompleteModule,
    MatChipsModule,
    MatButtonModule,
    CreateLayoutModule,
    InputModule,
    CardModule,
    MatCheckboxModule,
    ProjectRolesTableModule,
    MatIconModule,
    MatTooltipModule,
    HasRolePipeModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    FormsModule,
    TranslateModule,
    SearchOrgAutocompleteModule,
  ],
})
export default class ProjectGrantCreateModule {}
