import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { OwnedProjectDetailComponent } from './owned-project-detail.component';

const routes: Routes = [
  {
    path: '',
    component: OwnedProjectDetailComponent,
    data: {
      animation: 'HomePage',
      roles: ['project.read'],
    },
    canActivate: [RoleGuard],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class OwnedProjectDetailRoutingModule {}
