import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProjectRolesComponent } from './project-roles.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectRolesComponent,
    data: { animation: 'HomePage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectRolesRoutingModule {}
