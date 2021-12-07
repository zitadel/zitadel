import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';

import { GrantedProjectsComponent } from './granted-projects.component';

const routes: Routes = [
  {
    path: '',
    component: GrantedProjectsComponent,
    data: { animation: 'HomePage' },
  },
  {
    path: ':projectid/grant/:grantid/members',
    data: {
      type: ProjectType.PROJECTTYPE_GRANTED,
      roles: ['project.grant.member.read'],
    },
    loadChildren: () => import('src/app/modules/project-members/project-members.module').then((m) => m.ProjectMembersModule),
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class GrantedProjectsRoutingModule {}
