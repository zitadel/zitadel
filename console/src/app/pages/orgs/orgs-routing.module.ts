import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RoleGuard } from 'src/app/guards/role.guard';

import { OrgCreateComponent } from './org-create/org-create.component';
import { OrgDetailComponent } from './org-detail/org-detail.component';
import { OrgGridComponent } from './org-grid/org-grid.component';
import { PasswordPolicyComponent, PolicyComponentAction } from './password-policy/password-policy.component';

const routes: Routes = [
    {
        path: 'create',
        component: OrgCreateComponent,
        canActivate: [RoleGuard],
        data: {
            roles: ['iam.write'],
        },
        loadChildren: () => import('./org-create/org-create.module').then(m => m.OrgCreateModule),
    },
    {
        path: 'policy/:policytype/create',
        component: PasswordPolicyComponent,
        data: {
            action: PolicyComponentAction.CREATE,
        },
    },
    /// TODO: add roleguard for iam policy
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
