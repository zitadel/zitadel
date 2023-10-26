import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { ProjectRolesComponent } from './project-roles.component';

@NgModule({
  declarations: [ProjectRolesComponent],
  imports: [
    CommonModule,
    HasRoleModule,
    ProjectRolesTableModule,
    ReactiveFormsModule,
    HasRolePipeModule,
    InputModule,
    TranslateModule,
    MatButtonModule,
    FormsModule,
  ],
  exports: [ProjectRolesComponent],
})
export default class ProjectRolesModule {}
