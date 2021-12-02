import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProjectGrantsComponent } from './project-grants.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectGrantsComponent,
    data: { animation: 'HomePage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectGrantsRoutingModule {}
