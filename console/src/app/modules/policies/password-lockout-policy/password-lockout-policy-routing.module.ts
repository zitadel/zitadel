import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PolicyComponentAction } from '../policy-component-action.enum';
import { PasswordLockoutPolicyComponent } from './password-lockout-policy.component';

const routes: Routes = [
    {
        path: '',
        component: PasswordLockoutPolicyComponent,
        data: {
            animation: 'DetailPage',
            action: PolicyComponentAction.MODIFY,
        },
    },
    {
        path: 'create',
        component: PasswordLockoutPolicyComponent,
        data: {
            animation: 'DetailPage',
            action: PolicyComponentAction.CREATE,
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class PasswordLockoutPolicyRoutingModule { }
