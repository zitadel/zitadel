import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';

import { OrgDetailComponent } from './org-detail/org-detail.component';

const routes: Routes = [
  {
    path: 'members',
    loadChildren: () => import('./org-members/org-members.module'),
  },
  {
    path: 'provider',
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['org.idp.write'],
      serviceType: PolicyComponentServiceType.MGMT,
    },
    children: [
      {
        path: 'oauth',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-oauth/provider-oauth.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-oauth/provider-oauth.module'),
          },
        ],
      },
      {
        path: 'oidc',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-oidc/provider-oidc.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-oidc/provider-oidc.module'),
          },
        ],
      },
      {
        path: 'github-es',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-github-es/provider-github-es.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-github-es/provider-github-es.module'),
          },
        ],
      },
      {
        path: 'jwt',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-jwt/provider-jwt.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-jwt/provider-jwt.module'),
          },
        ],
      },
      {
        path: 'google',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-google/provider-google.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-google/provider-google.module'),
          },
        ],
      },
      {
        path: 'gitlab',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-gitlab/provider-gitlab.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-gitlab/provider-gitlab.module'),
          },
        ],
      },
      {
        path: 'gitlab-self-hosted',
        children: [
          {
            path: 'create',
            loadChildren: () =>
              import('src/app/modules/providers/provider-gitlab-self-hosted/provider-gitlab-self-hosted.module'),
          },
          {
            path: ':id',
            loadChildren: () =>
              import('src/app/modules/providers/provider-gitlab-self-hosted/provider-gitlab-self-hosted.module'),
          },
        ],
      },
      {
        path: 'github',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-github/provider-github.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-github/provider-github.module'),
          },
        ],
      },
      {
        path: 'azure-ad',
        children: [
          {
            path: 'create',
            loadChildren: () => import('src/app/modules/providers/provider-azure-ad/provider-azure-ad.module'),
          },
          {
            path: ':id',
            loadChildren: () => import('src/app/modules/providers/provider-azure-ad/provider-azure-ad.module'),
          },
        ],
      },
    ],
  },
  {
    path: '',
    component: OrgDetailComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class OrgRoutingModule {}
