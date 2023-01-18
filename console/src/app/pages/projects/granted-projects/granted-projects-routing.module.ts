import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';

import { GrantedProjectDetailComponent } from './granted-project-detail/granted-project-detail.component';

const routes: Routes = [
  {
    path: ':projectid/grant/:grantid/members',
    data: {
      type: ProjectType.PROJECTTYPE_GRANTED,
      roles: ['project.grant.member.read'],
    },
    loadChildren: () => import('src/app/modules/project-members/project-members.module'),
  },
  {
    path: ':id/grant/:grantId',
    component: GrantedProjectDetailComponent,
    data: { animation: 'HomePage' },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class GrantedProjectsRoutingModule {}
