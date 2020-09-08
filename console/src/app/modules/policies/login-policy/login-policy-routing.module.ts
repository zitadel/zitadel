import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PolicyComponentAction } from '../policy-component-action.enum';
import { LoginPolicyComponent } from './login-policy.component';

const routes: Routes = [
    {
        path: '',
        component: LoginPolicyComponent,
        data: {
            animation: 'DetailPage',
            action: PolicyComponentAction.MODIFY,
        },
    },
    {
        path: 'create',
        component: LoginPolicyComponent,
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
export class LoginPolicyRoutingModule { }
