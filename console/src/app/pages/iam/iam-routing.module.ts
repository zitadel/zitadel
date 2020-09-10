import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';
import { PolicyComponentServiceType, PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';

import { IamComponent } from './iam.component';

const routes: Routes = [
    {
        path: '',
        component: IamComponent,
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.read'],
        },
    },
    {
        path: 'members',
        loadChildren: () => import('./iam-members/iam-members.module').then(m => m.IamMembersModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.member.read'],
        },
    },
    {
        path: 'idp/create',
        loadChildren: () => import('src/app/modules/idp-create/idp-create.module').then(m => m.IdpCreateModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.idp.write'],
            serviceType: PolicyComponentServiceType.ADMIN,
        },
    },
    {
        path: 'idp/:id',
        loadChildren: () => import('src/app/modules/idp/idp.module').then(m => m.IdpModule),
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.idp.read'],
            serviceType: PolicyComponentServiceType.ADMIN,
        },
    },
    {
        path: `policy/${PolicyComponentType.LOGIN}`,
        data: {
            serviceType: PolicyComponentServiceType.ADMIN,
        },
        loadChildren: () => import('src/app/modules/policies/login-policy/login-policy.module')
            .then(m => m.LoginPolicyModule),
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class IamRoutingModule { }
