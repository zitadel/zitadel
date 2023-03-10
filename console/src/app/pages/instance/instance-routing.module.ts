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
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IamRoutingModule {}
