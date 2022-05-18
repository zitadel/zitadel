import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { QuicklinkStrategy } from 'ngx-quicklink';

import { AuthGuard } from './guards/auth.guard';
import { RoleGuard } from './guards/role.guard';
import { UserGrantContext } from './modules/user-grants/user-grants-datasource';
import { OrgCreateComponent } from './pages/org-create/org-create.component';

const routes: Routes = [
  {
    path: '',
    loadChildren: () => import('./pages/home/home.module').then((m) => m.HomeModule),
    canActivate: [AuthGuard],
  },
  {
    path: 'orgs',
    loadChildren: () => import('./pages/org-list/org-list.module').then((m) => m.OrgListModule),
    canActivate: [AuthGuard],
  },
  {
    path: 'granted-projects',
    loadChildren: () =>
      import('./pages/projects/granted-projects/granted-projects.module').then((m) => m.GrantedProjectsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['project.grant.read'],
    },
  },
  {
    path: 'projects',
    loadChildren: () => import('./pages/projects/projects.module').then((m) => m.ProjectsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['project.read'],
    },
  },
  {
    path: 'users',
    canActivate: [AuthGuard],
    children: [
      {
        path: '',
        loadChildren: () => import('src/app/pages/users/users.module').then((m) => m.UsersModule),
      },
    ],
  },
  {
    path: 'instance',
    loadChildren: () => import('./pages/instance/instance.module').then((m) => m.InstanceModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.read', 'iam.write'],
    },
  },
  {
    path: 'org/create',
    component: OrgCreateComponent,
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['(org.create)?(iam.write)?'],
    },
    loadChildren: () => import('./pages/org-create/org-create.module').then((m) => m.OrgCreateModule),
  },
  {
    path: 'org',
    loadChildren: () => import('./pages/orgs/org.module').then((m) => m.OrgModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['org.read'],
    },
  },
  {
    path: 'actions',
    loadChildren: () => import('./pages/actions/actions.module').then((m) => m.ActionsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['org.read'],
    },
  },
  {
    path: 'grants',
    loadChildren: () => import('./pages/grants/grants.module').then((m) => m.GrantsModule),
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
        loadChildren: () =>
          import('src/app/pages/user-grant-create/user-grant-create.module').then((m) => m.UserGrantCreateModule),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: 'project/:projectid',
        loadChildren: () =>
          import('src/app/pages/user-grant-create/user-grant-create.module').then((m) => m.UserGrantCreateModule),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: 'user/:userid',
        loadChildren: () =>
          import('src/app/pages/user-grant-create/user-grant-create.module').then((m) => m.UserGrantCreateModule),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
      {
        path: '',
        loadChildren: () =>
          import('src/app/pages/user-grant-create/user-grant-create.module').then((m) => m.UserGrantCreateModule),
        canActivate: [RoleGuard],
        data: {
          roles: ['user.grant.write'],
        },
      },
    ],
  },
  {
    path: 'failed-events',
    loadChildren: () => import('./pages/failed-events/failed-events.module').then((m) => m.FailedEventsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.read'],
    },
  },
  {
    path: 'views',
    loadChildren: () => import('./pages/iam-views/iam-views.module').then((m) => m.IamViewsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.read'],
    },
  },
  {
    path: 'settings',
    loadChildren: () => import('./pages/instance-settings/instance-settings.module').then((m) => m.InstanceSettingsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.read', 'iam.write'],
    },
  },
  {
    path: 'domains',
    loadChildren: () => import('./pages/domains/domains.module').then((m) => m.DomainsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['org.read'],
    },
  },
  {
    path: 'org-settings',
    loadChildren: () => import('./pages/org-settings/org-settings.module').then((m) => m.OrgSettingsModule),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['org.read', 'org.write'],
    },
  },
  {
    path: 'signedout',
    loadChildren: () => import('./pages/signedout/signedout.module').then((m) => m.SignedoutModule),
  },
  {
    path: '**',
    redirectTo: '/',
  },
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {
      preloadingStrategy: QuicklinkStrategy,
      relativeLinkResolution: 'legacy',
      scrollPositionRestoration: 'enabled',
    }),
  ],
  exports: [RouterModule],
})
export class AppRoutingModule {}
