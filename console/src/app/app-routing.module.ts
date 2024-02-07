import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthGuard } from './guards/auth.guard';
import { RoleGuard } from './guards/role.guard';
import { UserGrantContext } from './modules/user-grants/user-grants-datasource';
import { OrgCreateComponent } from './pages/org-create/org-create.component';

const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./pages/home/home.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['.'],
    },
  },
  {
    path: 'signedout',
    loadChildren: () => import('./pages/signedout/signedout.module'),
  },
  {
    path: 'orgs/create',
    component: OrgCreateComponent,
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['(org.create)?(iam.write)?'],
    },
    loadChildren: () => import('./pages/org-create/org-create.module'),
  },
  {
    path: 'orgs',
    loadChildren: () => import('./pages/org-list/org-list.module'),
    canActivate: [AuthGuard],
  },
  {
    path: 'granted-projects',
    loadChildren: () => import('./pages/projects/granted-projects/granted-projects.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['project.grant.read'],
    },
  },
  {
    path: 'projects',
    loadChildren: () => import('./pages/projects/projects.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['project.read'],
    },
  },
  {
    path: 'users',
    canActivate: [AuthGuard],
    loadChildren: () => import('src/app/pages/users/users.module'),
  },
  {
    path: 'instance',
    loadChildren: () => import('./pages/instance/instance.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.read', 'iam.write'],
    },
  },
  {
    path: 'org',
    loadChildren: () => import('./pages/orgs/org.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['org.read'],
    },
  },
  {
    path: 'actions',
    loadChildren: () => import('./pages/actions/actions.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['org.action.read', 'org.flow.read'],
    },
  },
  {
    path: 'grants',
    loadChildren: () => import('./pages/grants/grants.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      context: UserGrantContext.NONE,
      roles: ['user.grant.read'],
    },
  },
  {
    path: 'grant-create',
    canActivate: [AuthGuard],
    children: [
      {
        path: 'project/:projectid/grant/:grantid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: 'project/:projectid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: 'user/:userid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: '',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
    ],
  },
  {
    path: 'org-settings',
    loadChildren: () => import('./pages/org-settings/org-settings.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['policy.read'],
    },
  },
  {
    path: '**',
    redirectTo: '/',
  },
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {
      scrollPositionRestoration: 'enabled',
    }),
  ],
  exports: [RouterModule],
})
export class AppRoutingModule {}
