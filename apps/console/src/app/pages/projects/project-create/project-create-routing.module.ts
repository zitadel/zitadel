import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProjectCreateComponent } from './project-create.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectCreateComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectCreateRoutingModule {}
