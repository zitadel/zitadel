import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from 'src/app/guards/auth.guard';
import { roleGuard } from 'src/app/guards/role-guard';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';

import { OrgDetailComponent } from './org-detail/org-detail.component';

const routes: Routes = [
  {
    path: 'members',
    loadChildren: () => import('./org-members/org-members.module'),
  },
  {
    path: 'provider',
    canActivate: [authGuard, roleGuard],
    loadChildren: () => import('src/app/modules/providers/providers.module'),
    data: {
      roles: ['org.idp.read'],
      serviceType: PolicyComponentServiceType.MGMT,
    },
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
