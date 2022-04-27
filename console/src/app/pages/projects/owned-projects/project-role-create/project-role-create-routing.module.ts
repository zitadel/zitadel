import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProjectRoleCreateComponent } from './project-role-create.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectRoleCreateComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectRoleCreateRoutingModule {}
