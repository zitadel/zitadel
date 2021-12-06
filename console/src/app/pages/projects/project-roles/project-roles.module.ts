import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { ProjectRoleDetailComponent } from './project-role-detail/project-role-detail.component';
import { ProjectRolesRoutingModule } from './project-roles-routing.module';
import { ProjectRolesComponent } from './project-roles.component';

@NgModule({
  declarations: [ProjectRolesComponent, ProjectRoleDetailComponent],
  imports: [
    CommonModule,
    ProjectRolesRoutingModule,
    HasRoleModule,
    ProjectRolesTableModule,
    ReactiveFormsModule,
    HasRolePipeModule,
    InputModule,
    TranslateModule,
    FormsModule,
  ],
  exports: [ProjectRolesComponent],
})
export class ProjectRolesModule {}
