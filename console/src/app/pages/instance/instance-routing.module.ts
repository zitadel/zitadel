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
    loadChildren: () => import('src/app/modules/providers/providers.module'),
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
