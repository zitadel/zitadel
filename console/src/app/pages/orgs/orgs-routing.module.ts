import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { OrgCreateComponent } from './org-create/org-create.component';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgGridComponent } from './org-grid/org-grid.component';

export enum PolicyComponentAction {
    CREATE = 'create',
    MODIFY = 'modify',
}

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
        path: 'policy/age',
        data: {
            action: PolicyComponentAction.MODIFY,
        },
        loadChildren: () => import('./password-age-policy/password-age-policy.module')
            .then(m => m.PasswordAgePolicyModule),
    },
    {
        path: 'policy/lockout/create',
        data: {
            action: PolicyComponentAction.CREATE,
        },
        loadChildren: () => import('./password-lockout-policy/password-lockout-policy.module')
            .then(m => m.PasswordLockoutPolicyModule),
    },
    {
        path: 'policy/lockout',
        data: {
            action: PolicyComponentAction.MODIFY,
        },
        loadChildren: () => import('./password-lockout-policy/password-lockout-policy.module')
            .then(m => m.PasswordLockoutPolicyModule),
    },
    {
        path: 'policy/complexity/create',
        data: {
            action: PolicyComponentAction.CREATE,
        },
        loadChildren: () => import('./password-complexity-policy/password-complexity-policy.module')
            .then(m => m.PasswordComplexityPolicyModule),
    },
    {
        path: 'policy/complexity',
        data: {
            action: PolicyComponentAction.MODIFY,
        },
        loadChildren: () => import('./password-complexity-policy/password-complexity-policy.module')
            .then(m => m.PasswordComplexityPolicyModule),
    },
    {
        path: 'policy/iam_policy/create',
        data: {
            action: PolicyComponentAction.CREATE,
        },
        loadChildren: () => import('./password-iam-policy/password-iam-policy.module')
            .then(m => m.PasswordIamPolicyModule),
    },
    {
        path: 'policy/iam_policy',
        data: {
            action: PolicyComponentAction.MODIFY,
        },
        loadChildren: () => import('./password-iam-policy/password-iam-policy.module')
            .then(m => m.PasswordIamPolicyModule),
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
