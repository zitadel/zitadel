import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProjectGrantCreateComponent } from './project-grant-create.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectGrantCreateComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectGrantCreateRoutingModule {}
