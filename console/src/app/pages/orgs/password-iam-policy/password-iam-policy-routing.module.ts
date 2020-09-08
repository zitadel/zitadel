import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { PasswordIamPolicyComponent } from './password-iam-policy.component';

const routes: Routes = [
    {
        path: '',
        component: PasswordIamPolicyComponent,
        data: { animation: 'DetailPage' },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class PasswordIamPolicyRoutingModule { }
