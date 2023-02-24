import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyAutocompleteModule as MatAutocompleteModule } from '@angular/material/legacy-autocomplete';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CardModule } from 'src/app/modules/card/card.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

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
  ],
})
export default class ProjectGrantCreateModule {}
