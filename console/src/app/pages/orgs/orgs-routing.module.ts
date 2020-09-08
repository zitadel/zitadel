import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { OrgCreateComponent } from './org-create/org-create.component';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgGridComponent } from './org-grid/org-grid.component';

const routes: Routes = [
    {
        path: 'create',
        component: OrgCreateComponent,
        canActivate: [RoleGuard],
        data: {
            roles: ['(org.create)?(iam.write)?'],
        },
        loadChildren: () => import('./org-create/org-create.module').then(m => m.OrgCreateModule),
    },
    {
        path: 'idp/create',
        loadChildren: () => import('src/app/modules/idp-create/idp-create.module').then(m => m.IdpCreateModule),
        canActivate: [RoleGuard],
        data: {
            roles: ['org.idp.write'],
        },
    },
    {
        path: 'policy',
        children: [
            {
                path: 'age',
                loadChildren: () => import('./policies/password-age-policy/password-age-policy.module')
                    .then(m => m.PasswordAgePolicyModule),
            },
            {
                path: 'lockout',
                loadChildren: () => import('./policies/password-lockout-policy/password-lockout-policy.module')
                    .then(m => m.PasswordLockoutPolicyModule),
            },
            {
                path: 'complexity',
                loadChildren: () => import('./policies/password-complexity-policy/password-complexity-policy.module')
                    .then(m => m.PasswordComplexityPolicyModule),
            },
            {
                path: 'iam_policy',
                loadChildren: () => import('./policies/password-iam-policy/password-iam-policy.module')
                    .then(m => m.PasswordIamPolicyModule),
            },
        ],
    },
    {
        path: 'members',
        loadChildren: () => import('./org-members/org-members.module').then(m => m.OrgMembersModule),
    },
    {
        path: '',
        component: OrgDetailComponent,
    },
    {
        path: 'overview',
        component: OrgGridComponent,
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class OrgsRoutingModule { }
