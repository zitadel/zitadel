import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from 'src/app/guards/auth.guard';
import { roleGuard } from 'src/app/guards/role-guard';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';

import { InstanceComponent } from './instance.component';

const routes: Routes = [
  {
    path: '',
    component: InstanceComponent,
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['iam.read'],
    },
    children: [
      {
        path: 'actions',
        loadChildren: () => import('src/app/modules/actions-two/actions-two.module'),
      },
    ],
  },
  {
    path: 'members',
    loadChildren: () => import('./instance-members/instance-members.module'),
    canActivate: [authGuard, roleGuard],
    data: {
      roles: ['iam.member.read'],
    },
  },
  {
    path: 'provider',
    canActivate: [authGuard, roleGuard],
    loadChildren: () => import('src/app/modules/providers/providers.module'),
    data: {
      roles: ['iam.idp.read'],
      serviceType: PolicyComponentServiceType.ADMIN,
    },
  },
  {
    path: 'smtpprovider',
    canActivate: [authGuard, roleGuard],
    loadChildren: () => import('src/app/modules/smtp-provider/smtp-provider.module'),
    data: {
      roles: ['iam.idp.read'],
      serviceType: PolicyComponentServiceType.ADMIN,
    },
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class IamRoutingModule {}
