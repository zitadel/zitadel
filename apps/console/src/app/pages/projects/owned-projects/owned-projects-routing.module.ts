import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { roleGuard } from 'src/app/guards/role-guard';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';

const routes: Routes = [
  {
    path: '',
    data: {
      animation: 'HomePage',
      roles: ['project.read'],
    },
    canActivate: [roleGuard],
    loadChildren: () => import('./owned-project-detail/owned-project-detail.module'),
  },
  {
    path: 'members',
    data: {
      type: ProjectType.PROJECTTYPE_OWNED,
      roles: ['project.member.read'],
    },
    canActivate: [roleGuard],
    loadChildren: () => import('src/app/modules/project-members/project-members.module'),
  },
  {
    path: 'apps',
    data: {
      animation: 'AddPage',
      roles: ['project.app.read'],
    },
    canActivate: [roleGuard],
    loadChildren: () => import('src/app/pages/projects/apps/apps.module'),
  },
  {
    path: 'projectgrants',
    data: {
      // animation: 'AddPage',
      // roles: ['project.grant.read:' + ':projectid', 'project.grant.read'],
    },
    // canActivate: [RoleGuard],
    loadChildren: () => import('src/app/pages/projects/owned-projects/project-grants/project-grants.module'),
  },
  {
    path: 'roles',
    loadChildren: () => import('src/app/pages/projects/owned-projects/project-roles/project-roles.module'),
  },
  {
    path: 'roles/create',
    loadChildren: () => import('./project-role-create/project-role-create.module'),
  },
  {
    path: 'projectgrants/create',
    loadChildren: () => import('./project-grant-create/project-grant-create.module'),
  },
  {
    path: 'projectgrants/:grantid',
    loadChildren: () => import('./project-grant-detail/project-grant-detail.module'),
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class OwnedProjectsRoutingModule {}
