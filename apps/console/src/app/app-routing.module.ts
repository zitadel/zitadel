import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { authGuard } from './guards/auth.guard';
import { roleGuard } from './guards/role-guard';
import { UserGrantContext } from './modules/user-grants/user-grants-datasource';
import { OrgCreateComponent } from './pages/org-create/org-create.component';

const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./pages/home/home.module'),
    canActivate: [authGuard, roleGuard],
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
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['(org.create)?(iam.write)?'],
    },
    loadChildren: () => import('./pages/org-create/org-create.module'),
  },
  {
    path: 'orgs',
    loadChildren: () => import('./pages/org-list/org-list.module'),
    canActivate: [authGuard],
  },
  {
    path: 'granted-projects',
    loadChildren: () => import('./pages/projects/granted-projects/granted-projects.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['project.grant.read'],
    },
  },
  {
    path: 'projects',
    loadChildren: () => import('./pages/projects/projects.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['project.read'],
    },
  },
  {
    path: 'users',
    canActivate: [authGuard],
    loadChildren: () => import('src/app/pages/users/users.module'),
  },
  {
    path: 'instance',
    loadChildren: () => import('./pages/instance/instance.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['iam.read', 'iam.write'],
    },
  },
  {
    path: 'org',
    loadChildren: () => import('./pages/orgs/org.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['org.read'],
    },
  },
  {
    path: 'actions',
    loadChildren: () => import('./pages/actions/actions.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['org.action.read', 'org.flow.read'],
    },
  },
  {
    path: 'grants',
    loadChildren: () => import('./pages/grants/grants.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      context: UserGrantContext.NONE,
      roles: ['user.grant.read'],
    },
  },
  {
    path: 'grant-create',
    canActivate: [authGuard],
    children: [
      {
        path: 'project/:projectid/grant/:grantid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [roleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: 'project/:projectid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [roleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: 'user/:userid',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [roleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: '',
        loadChildren: () => import('src/app/pages/user-grant-create/user-grant-create.module'),
        canActivate: [roleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
    ],
  },
  {
    path: 'org-settings',
    loadChildren: () => import('./pages/org-settings/org-settings.module'),
    canActivate: [authGuard, roleGuard],
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
