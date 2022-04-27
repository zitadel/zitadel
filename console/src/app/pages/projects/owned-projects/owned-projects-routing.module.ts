import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';

const routes: Routes = [
  {
    path: '',
    data: {
      animation: 'HomePage',
      roles: ['project.read'],
    },
    canActivate: [RoleGuard],
    loadChildren: () => import('./owned-project-detail/owned-project-detail.module').then((m) => m.OwnedProjectDetailModule),
  },
  {
    path: 'members',
    data: {
      type: ProjectType.PROJECTTYPE_OWNED,
      roles: ['project.member.read'],
    },
    canActivate: [RoleGuard],
    loadChildren: () => import('src/app/modules/project-members/project-members.module').then((m) => m.ProjectMembersModule),
  },
  {
    path: 'apps',
    data: {
      animation: 'AddPage',
      roles: ['project.app.read'],
    },
    canActivate: [RoleGuard],
    loadChildren: () => import('src/app/pages/projects/apps/apps.module').then((m) => m.AppsModule),
  },
  {
    path: 'projectgrants',
    data: {
      // animation: 'AddPage',
      // roles: ['project.grant.read:' + ':projectid', 'project.grant.read'],
    },
    // canActivate: [RoleGuard],
    loadChildren: () =>
      import('src/app/pages/projects/owned-projects/project-grants/project-grants.module').then(
        (m) => m.ProjectGrantsModule,
      ),
  },
  {
    path: 'grants',
    loadChildren: () => import('src/app/pages/grants/grants.module').then((m) => m.GrantsModule),
    canActivate: [RoleGuard],
    data: {
      roles: ['user.grant.read'],
      context: UserGrantContext.OWNED_PROJECT,
    },
  },
  {
    path: 'roles',
    loadChildren: () =>
      import('src/app/pages/projects/owned-projects/project-roles/project-roles.module').then((m) => m.ProjectRolesModule),
  },
  {
    path: 'roles/create',
    loadChildren: () => import('./project-role-create/project-role-create.module').then((m) => m.ProjectRoleCreateModule),
  },
  {
    path: 'projectgrants/create',
    loadChildren: () => import('./project-grant-create/project-grant-create.module').then((m) => m.ProjectGrantCreateModule),
  },
  {
    path: 'projectgrants/:grantid',
    loadChildren: () => import('./project-grant-detail/project-grant-detail.module').then((m) => m.ProjectGrantDetailModule),
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class OwnedProjectsRoutingModule {}
