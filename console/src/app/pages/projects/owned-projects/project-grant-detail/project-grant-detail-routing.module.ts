import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProjectGrantDetailComponent } from './project-grant-detail.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectGrantDetailComponent,
    data: { animation: 'AddPage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectGrantDetailRoutingModule {}
