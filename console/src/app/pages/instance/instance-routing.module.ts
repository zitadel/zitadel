import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';

import { InstanceComponent } from './instance.component';

const routes: Routes = [
  {
    path: '',
    component: InstanceComponent,
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.read'],
    },
  },
  {
    path: 'members',
    loadChildren: () => import('./instance-members/instance-members.module'),
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.member.read'],
    },
  },
  {
    path: 'provider',
    canActivate: [AuthGuard, RoleGuard],
    data: {
      roles: ['iam.idp.write'],
      serviceType: PolicyComponentServiceType.ADMIN,
    },
    children: [
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
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IamRoutingModule {}
