import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { RoleGuard } from 'src/app/guards/role.guard';
import { FeatureServiceType } from 'src/app/modules/features/features.component';
import { PolicyComponentServiceType, PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';

import { EventstoreComponent } from './eventstore/eventstore.component';
import { IamComponent } from './iam.component';

const routes: Routes = [
    {
        path: 'policies',
        component: IamComponent,
        canActivate: [AuthGuard, RoleGuard],
        data: {
            roles: ['iam.read'],
        },
    },
    {
        path: 'eventstore',
        component: EventstoreComponent,
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
        path: 'features',
        loadChildren: () => import('src/app/modules/features/features.module').then(m => m.FeaturesModule),
        // canActivate: [RoleGuard],
        data: {
            roles: ['iam.features.read'],
            serviceType: FeatureServiceType.ADMIN,
        },
    },
    {
        path: 'idp',
        children: [
            {
                path: 'create',
                loadChildren: () => import('src/app/modules/idp-create/idp-create.module').then(m => m.IdpCreateModule),
                canActivate: [AuthGuard, RoleGuard],
                data: {
                    roles: ['iam.idp.write'],
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
            },
            {
                path: ':id',
                loadChildren: () => import('src/app/modules/idp/idp.module').then(m => m.IdpModule),
                canActivate: [AuthGuard, RoleGuard],
                data: {
                    roles: ['iam.idp.read'],
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
            },
        ],
    },
    {
        path: 'policy',
        children: [
            {
                path: PolicyComponentType.AGE,
                data: {
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
                loadChildren: () => import('src/app/modules/policies/password-age-policy/password-age-policy.module')
                    .then(m => m.PasswordAgePolicyModule),
            },
            {
                path: PolicyComponentType.LOCKOUT,
                data: {
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
                loadChildren: () => import('src/app/modules/policies/password-lockout-policy/password-lockout-policy.module')
                    .then(m => m.PasswordLockoutPolicyModule),
            },
            {
                path: PolicyComponentType.COMPLEXITY,
                data: {
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
                loadChildren: () => import('src/app/modules/policies/password-complexity-policy/password-complexity-policy.module')
                    .then(m => m.PasswordComplexityPolicyModule),
            },
            {
                path: PolicyComponentType.IAM,
                data: {
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
                loadChildren: () => import('src/app/modules/policies/org-iam-policy/org-iam-policy.module')
                    .then(m => m.OrgIamPolicyModule),
            },
            {
                path: PolicyComponentType.LOGIN,
                data: {
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
                loadChildren: () => import('src/app/modules/policies/login-policy/login-policy.module')
                    .then(m => m.LoginPolicyModule),
            },
            {
                path: PolicyComponentType.LABEL,
                data: {
                    serviceType: PolicyComponentServiceType.ADMIN,
                },
                loadChildren: () => import('src/app/modules/policies/label-policy/label-policy.module')
                    .then(m => m.LabelPolicyModule),
            },
        ],
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class IamRoutingModule { }
