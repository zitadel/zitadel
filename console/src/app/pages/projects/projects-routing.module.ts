import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { ProjectsComponent } from './projects.component';

const routes: Routes = [
  {
    path: '',
    component: ProjectsComponent,
    data: { animation: 'HomePage' },
  },
  {
    path: 'create',
    loadChildren: () => import('./project-create/project-create.module').then((m) => m.ProjectCreateModule),
    canActivate: [RoleGuard],
    data: {
      animation: 'AddPage',
      roles: ['project.create'],
    },
  },
  {
    path: ':projectid',
    loadChildren: () => import('./owned-projects/owned-projects.module').then((m) => m.OwnedProjectsModule),
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ProjectsRoutingModule {}
