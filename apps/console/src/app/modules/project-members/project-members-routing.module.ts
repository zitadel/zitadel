import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProjectMembersComponent } from './project-members.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectMembersComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectMembersRoutingModule {}
