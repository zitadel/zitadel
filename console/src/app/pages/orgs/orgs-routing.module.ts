import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { OrgCreateComponent } from './org-create/org-create.component';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgGridComponent } from './org-grid/org-grid.component';
import { PasswordAgePolicyComponent } from './password-age-policy/password-age-policy.component';
import { PasswordLockoutPolicyComponent } from './password-lockout-policy/password-lockout-policy.component';
import { PasswordPolicyComponent, PolicyComponentAction } from './password-policy/password-policy.component';

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
        component: PasswordAgePolicyComponent,
        data: {
            action: PolicyComponentAction.MODIFY,
        },
        loadChildren: () => import('./password-age-policy/password-age-policy.module')
            .then(m => m.PasswordAgePolicyModule),
    },
    {
        path: 'policy/lockout',
        component: PasswordLockoutPolicyComponent,
        data: {
            action: PolicyComponentAction.MODIFY,
        },
        loadChildren: () => import('./password-lockout-policy/password-lockout-policy.module')
            .then(m => m.PasswordLockoutPolicyModule),
    },
    {
        path: 'policy/:policytype/create',
        component: PasswordPolicyComponent,
        data: {
            action: PolicyComponentAction.CREATE,
        },
    },
    {
        path: 'policy/:policytype',
        component: PasswordPolicyComponent,
        data: {
            action: PolicyComponentAction.MODIFY,
        },
        loadChildren: () => import('./password-policy/password-policy.module').then(m => m.PasswordPolicyModule),
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
