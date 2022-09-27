import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';

import { OrgDetailComponent } from './org-detail/org-detail.component';

const routes: Routes = [
  {
    path: 'idp',
    children: [
      {
        path: 'create',
        loadChildren: () => import('src/app/modules/idp-create/idp-create.module').then((m) => m.IdpCreateModule),
        canActivate: [RoleGuard],
        data: {
          roles: ['org.idp.write'],
          serviceType: PolicyComponentServiceType.MGMT,
        },
      },
      {
        path: ':id',
        loadChildren: () => import('src/app/modules/idp/idp.module').then((m) => m.IdpModule),
        canActivate: [RoleGuard],
        data: {
          roles: ['org.idp.read'],
          serviceType: PolicyComponentServiceType.MGMT,
        },
      },
    ],
  },
  {
    path: 'members',
    loadChildren: () => import('./org-members/org-members.module').then((m) => m.OrgMembersModule),
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
